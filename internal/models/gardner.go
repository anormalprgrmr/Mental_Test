package models

// GardnerTest returns a compact Gardner Multiple Intelligences style assessment.
func GardnerTest() TestDefinition {
	return TestDefinition{
		ID:          "gardner",
		Title:       "آزمون هوش‌های چندگانه گاردنر",
		Description: "سبک‌های هوشی غالب و توانمندی‌های فعلی خود را بهتر بشناسید.",
		Dimensions: []Dimension{
			{ID: "linguistic", Name: "هوش زبانی", Description: "کلمات، مطالعه، نوشتن و توضیح ایده‌ها."},
			{ID: "logical", Name: "هوش منطقی-ریاضی", Description: "الگوها، استدلال، اعداد و سیستم‌ها."},
			{ID: "spatial", Name: "هوش فضایی", Description: "تصویرسازی ذهنی، طراحی، نقشه و تجسم."},
			{ID: "bodily", Name: "هوش بدنی-جنبشی", Description: "حرکت، کار عملی و یادگیری با انجام دادن."},
			{ID: "musical", Name: "هوش موسیقایی", Description: "ریتم، صدا، ملودی و گوش دادن دقیق."},
			{ID: "interpersonal", Name: "هوش میان‌فردی", Description: "درک دیگران، همکاری و ارتباط با افراد."},
			{ID: "intrapersonal", Name: "هوش درون‌فردی", Description: "خودشناسی، هدف‌گذاری و آگاهی از احساسات."},
			{ID: "naturalistic", Name: "هوش طبیعت‌گرا", Description: "طبیعت، موجودات زنده، مشاهده و دسته‌بندی."},
		},
		Questions: []Question{
			likertQuestion("از خواندن، نوشتن یا تعریف کردن داستان لذت می‌برم.", "linguistic"),
			likertQuestion("حل کردن معما، محاسبه یا مسئله‌های منطقی را دوست دارم.", "logical"),
			likertQuestion("اطلاعات را با نمودار، رنگ یا نقشه بهتر به خاطر می‌سپارم.", "spatial"),
			likertQuestion("با تمرین عملی یا ساختن چیزی بهتر یاد می‌گیرم.", "bodily"),
			likertQuestion("موسیقی، ریتم یا الگوهای صوتی به تمرکز یا حافظه من کمک می‌کنند.", "musical"),
			likertQuestion("دیگران معمولاً برای مشورت یا کار گروهی به من مراجعه می‌کنند.", "interpersonal"),
			likertQuestion("زمانی را صرف فکر کردن درباره هدف‌ها، ارزش‌ها و احساساتم می‌کنم.", "intrapersonal"),
			likertQuestion("در طبیعت، حیوانات، گیاهان یا محیط اطراف الگوها را سریع متوجه می‌شوم.", "naturalistic"),
			likertQuestion("می‌توانم ایده‌های پیچیده را با کلمات ساده و واضح توضیح بدهم.", "linguistic"),
			likertQuestion("دوست دارم قانون‌ها و سازوکار پشت اتفاقات را پیدا کنم.", "logical"),
			likertQuestion("می‌توانم اشیا یا فضاها را به‌خوبی در ذهنم تصور کنم.", "spatial"),
			likertQuestion("ورزش، کاردستی یا فعالیت‌های حرکتی به من انرژی می‌دهند.", "bodily"),
			likertQuestion("خیلی زود متوجه می‌شوم موسیقی از ریتم یا کوک خارج شده است.", "musical"),
			likertQuestion("حال‌وهوا و انگیزه‌های دیگران را معمولاً خوب درک می‌کنم.", "interpersonal"),
			likertQuestion("قبل از تصمیم‌گیری ترجیح می‌دهم در سکوت خودم را بهتر بررسی کنم.", "intrapersonal"),
			likertQuestion("از مشاهده، دسته‌بندی یا مراقبت از چیزهای طبیعی لذت می‌برم.", "naturalistic"),
		},
	}
}

func likertQuestion(text, dimensionID string) Question {
	return Question{
		Text: text,
		Options: []Option{
			{Text: "کاملاً مخالفم", Scores: map[string]int{dimensionID: 1}},
			{Text: "مخالفم", Scores: map[string]int{dimensionID: 2}},
			{Text: "نظری ندارم", Scores: map[string]int{dimensionID: 3}},
			{Text: "موافقم", Scores: map[string]int{dimensionID: 4}},
			{Text: "کاملاً موافقم", Scores: map[string]int{dimensionID: 5}},
		},
	}
}
