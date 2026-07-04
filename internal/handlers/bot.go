package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"mental_test/internal/models"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot     *tele.Bot
	store   ResultStore
	catalog *models.Catalog
	admins  []int64
	logger  *slog.Logger

	startTestButton tele.Btn
	answerButton    tele.Btn

	mu       sync.RWMutex
	sessions map[int64]*Session
}

type Session struct {
	TestID  string
	Answers []int
}

// ResultStore is the persistence dependency required by the Telegram handlers.
// Keeping it as an interface prevents handlers from depending on a specific DB implementation.
type ResultStore interface {
	UpsertUser(ctx context.Context, user models.TelegramUser) error
	SaveResult(ctx context.Context, userID int64, testID string, result models.Result) (models.SavedResult, error)
}

type Config struct {
	Token  string
	Poller time.Duration
	Admins []int64
}

func New(cfg Config, store ResultStore, catalog *models.Catalog, logger *slog.Logger) (*Bot, error) {
	teleBot, err := tele.NewBot(tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: cfg.Poller},
	})
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	if logger == nil {
		logger = slog.Default()
	}

	b := &Bot{
		bot:             teleBot,
		store:           store,
		catalog:         catalog,
		admins:          append([]int64(nil), cfg.Admins...),
		logger:          logger,
		startTestButton: tele.Btn{Unique: "start_test"},
		answerButton:    tele.Btn{Unique: "answer"},
		sessions:        make(map[int64]*Session),
	}
	b.registerHandlers()
	return b, nil
}

func (b *Bot) Start() {
	b.logger.Info("telegram bot started")
	b.bot.Start()
}

func (b *Bot) registerHandlers() {
	b.bot.Handle("/start", b.handleStart)
	b.bot.Handle("/tests", b.handleTests)
	b.bot.Handle("/cancel", b.handleCancel)
	b.bot.Handle(&b.startTestButton, b.handleStartTest)
	b.bot.Handle(&b.answerButton, b.handleAnswer)
}

func (b *Bot) handleStart(c tele.Context) error {
	if err := b.saveUser(c); err != nil {
		return c.Send("امکان ذخیره اطلاعات شما وجود نداشت. لطفاً کمی بعد دوباره تلاش کنید.")
	}

	return c.Send("سلام! خوش آمدید. برای شروع، یکی از آزمون‌ها را انتخاب کنید. هر زمان هم می‌توانید از دستور /tests استفاده کنید.", b.testsMarkup())
}

func (b *Bot) handleTests(c tele.Context) error {
	if err := b.saveUser(c); err != nil {
		return c.Send("امکان ذخیره اطلاعات شما وجود نداشت. لطفاً کمی بعد دوباره تلاش کنید.")
	}
	return c.Send("آزمون‌های موجود:", b.testsMarkup())
}

func (b *Bot) handleCancel(c tele.Context) error {
	userID := c.Sender().ID
	b.mu.Lock()
	delete(b.sessions, userID)
	b.mu.Unlock()
	return c.Send("آزمون فعلی لغو شد. برای شروع دوباره از /tests استفاده کنید.")
}

func (b *Bot) handleStartTest(c tele.Context) error {
	if err := c.Respond(); err != nil {
		b.logger.Warn("respond to start test callback", "error", err)
	}
	if err := b.saveUser(c); err != nil {
		return c.Send("امکان ذخیره اطلاعات شما وجود نداشت. لطفاً کمی بعد دوباره تلاش کنید.")
	}

	testID := c.Data()
	test, ok := b.catalog.Get(testID)
	if !ok {
		return c.Send("این آزمون شناخته نشد. لطفاً از /tests یکی از آزمون‌های موجود را انتخاب کنید.")
	}

	userID := c.Sender().ID
	b.mu.Lock()
	b.sessions[userID] = &Session{TestID: testID, Answers: make([]int, 0, len(test.Questions))}
	b.mu.Unlock()

	return b.sendQuestion(c, test, 0)
}

func (b *Bot) handleAnswer(c tele.Context) error {
	if err := c.Respond(); err != nil {
		b.logger.Warn("respond to answer callback", "error", err)
	}

	userID := c.Sender().ID
	b.mu.RLock()
	session, ok := b.sessions[userID]
	b.mu.RUnlock()
	if !ok {
		return c.Send("در حال حاضر آزمون فعالی ندارید. برای شروع از /tests استفاده کنید.")
	}

	test, ok := b.catalog.Get(session.TestID)
	if !ok {
		b.deleteSession(userID)
		return c.Send("این آزمون دیگر در دسترس نیست. لطفاً با /tests آزمون دیگری انتخاب کنید.")
	}

	answerIndex, err := parseAnswerData(c.Data())
	if err != nil {
		return c.Send("پاسخ نامعتبر است. لطفاً فقط از دکمه‌های پاسخ استفاده کنید.")
	}

	b.mu.Lock()
	session = b.sessions[userID]
	questionIndex := len(session.Answers)
	if questionIndex >= len(test.Questions) {
		b.mu.Unlock()
		return c.Send("این آزمون قبلاً کامل شده است. برای شروع آزمون دیگر از /tests استفاده کنید.")
	}
	if answerIndex < 0 || answerIndex >= len(test.Questions[questionIndex].Options) {
		b.mu.Unlock()
		return c.Send("پاسخ انتخاب‌شده برای این سؤال نامعتبر است.")
	}
	session.Answers = append(session.Answers, answerIndex)
	answers := append([]int(nil), session.Answers...)
	completed := len(answers) == len(test.Questions)
	b.mu.Unlock()

	if !completed {
		return b.sendQuestion(c, test, len(answers))
	}

	b.deleteSession(userID)
	result, err := test.Calculate(answers)
	if err != nil {
		b.logger.Error("calculate result", "user_id", userID, "test_id", test.ID, "error", err)
		return c.Send("امکان محاسبه نتیجه شما وجود نداشت. لطفاً کمی بعد دوباره تلاش کنید.")
	}

	saved, err := b.store.SaveResult(context.Background(), userID, test.ID, result)
	if err != nil {
		b.logger.Error("save result", "user_id", userID, "test_id", test.ID, "error", err)
		return c.Send("امکان ذخیره نتیجه شما وجود نداشت. لطفاً کمی بعد دوباره تلاش کنید.")
	}

	message := formatUserResult(result)
	if err := c.Send(message, tele.ModeMarkdown); err != nil {
		return err
	}

	b.notifyAdmins(c.Sender(), saved)
	return nil
}

func (b *Bot) testsMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0)

	for _, test := range b.catalog.All() {
		button := b.startTestButton
		button.Text = test.Title
		button.Data = test.ID
		rows = append(rows, markup.Row(button))
	}
	markup.Inline(rows...)
	return markup
}

func (b *Bot) sendQuestion(c tele.Context, test models.TestDefinition, index int) error {
	question := test.Questions[index]
	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0, len(question.Options))

	for optionIndex, option := range question.Options {
		button := b.answerButton
		button.Text = option.Text
		button.Data = fmt.Sprintf("%d", optionIndex)
		rows = append(rows, markup.Row(button))
	}
	markup.Inline(rows...)

	text := fmt.Sprintf("%s\n\nسؤال %d از %d:\n%s", test.Title, index+1, len(test.Questions), question.Text)
	return c.Send(text, markup)
}

func (b *Bot) saveUser(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return nil
	}

	return b.store.UpsertUser(context.Background(), models.TelegramUser{
		ID:        sender.ID,
		Username:  sender.Username,
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
	})
}

func (b *Bot) notifyAdmins(user *tele.User, result models.SavedResult) {
	if len(b.admins) == 0 {
		return
	}

	message := formatAdminResult(user, result)
	for _, adminID := range b.admins {
		recipient := &tele.User{ID: adminID}
		if _, err := b.bot.Send(recipient, message, tele.ModeMarkdown); err != nil {
			b.logger.Warn("notify admin failed", "admin_id", adminID, "result_id", result.ID, "error", err)
		}
	}
}

func (b *Bot) deleteSession(userID int64) {
	b.mu.Lock()
	delete(b.sessions, userID)
	b.mu.Unlock()
}

func parseAnswerData(data string) (int, error) {
	var answer int
	_, err := fmt.Sscanf(strings.TrimSpace(data), "%d", &answer)
	return answer, err
}

func formatUserResult(result models.Result) string {
	var builder strings.Builder
	builder.WriteString("✅ *نتیجه آزمون شما آماده است!*\n\n")
	builder.WriteString(escapeMarkdown(result.Summary))
	builder.WriteString("\n\n*امتیازها:*\n")
	for _, score := range result.Scores {
		builder.WriteString(fmt.Sprintf("• %s: %d/%d\n", escapeMarkdown(score.Dimension.Name), score.Score, score.MaxPossible))
	}
	builder.WriteString("\nبرای انجام آزمون دیگر از /tests استفاده کنید.")
	return builder.String()
}

func formatAdminResult(user *tele.User, saved models.SavedResult) string {
	name := "کاربر ناشناس"
	username := ""
	if user != nil {
		name = strings.TrimSpace(strings.Join([]string{user.FirstName, user.LastName}, " "))
		if name == "" {
			name = fmt.Sprintf("شناسه %d", user.ID)
		}
		if user.Username != "" {
			username = "@" + user.Username
		}
	}

	var builder strings.Builder
	builder.WriteString("📌 *نتیجه آزمون جدید*\n\n")
	builder.WriteString(fmt.Sprintf("شناسه نتیجه: `%d`\n", saved.ID))
	builder.WriteString(fmt.Sprintf("کاربر: %s %s\n", escapeMarkdown(name), escapeMarkdown(username)))
	builder.WriteString(fmt.Sprintf("شناسه تلگرام: `%d`\n", saved.UserID))
	builder.WriteString(fmt.Sprintf("آزمون: %s\n\n", escapeMarkdown(saved.Result.Title)))
	for _, score := range saved.Result.Scores {
		builder.WriteString(fmt.Sprintf("• %s: %d/%d\n", escapeMarkdown(score.Dimension.Name), score.Score, score.MaxPossible))
	}
	return builder.String()
}

func escapeMarkdown(value string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
	)
	return replacer.Replace(value)
}
