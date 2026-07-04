package models

// CliftonTest returns a compact strengths-style assessment inspired by CliftonStrengths themes.
func CliftonTest() TestDefinition {
	return TestDefinition{
		ID:          "clifton",
		Title:       "آزمون سبک توانمندی‌های کلیفتون",
		Description: "تم‌های غالب کاری، انگیزشی و رفتاری خود را شناسایی کنید.",
		Dimensions: []Dimension{
			{ID: "achiever", Name: "دستاوردگرا", Description: "انگیزه، پشتکار و رضایت از انجام دادن کارها."},
			{ID: "strategic", Name: "استراتژیک", Description: "دیدن الگوها، مسیرهای جایگزین و راه‌حل‌های بهتر."},
			{ID: "learner", Name: "یادگیرنده", Description: "انرژی گرفتن از یادگیری، رشد و تسلط بر موضوعات جدید."},
			{ID: "relator", Name: "رابطه‌محور", Description: "اعتماد، عمق رابطه و ارتباط‌های نزدیک و پایدار."},
			{ID: "communication", Name: "ارتباط‌گر", Description: "تبدیل ایده‌ها به پیام‌های روشن، جذاب و قابل فهم."},
			{ID: "responsibility", Name: "مسئولیت‌پذیر", Description: "مالکیت کار، قابل اعتماد بودن و پایبندی به تعهدات."},
		},
		Questions: []Question{
			likertQuestion("وقتی در یک روز کارهای زیادی را کامل می‌کنم احساس رضایت دارم.", "achiever"),
			likertQuestion("برای حل یک مسئله، معمولاً سریع چند مسیر جایگزین می‌بینم.", "strategic"),
			likertQuestion("یاد گرفتن چیزهای جدید حتی قبل از استفاده عملی برایم هیجان‌انگیز است.", "learner"),
			likertQuestion("چند رابطه عمیق را به ارتباط‌های سطحی زیاد ترجیح می‌دهم.", "relator"),
			likertQuestion("از ارائه ایده‌ها، داستان‌ها یا توضیح دادن موضوعات به دیگران لذت می‌برم.", "communication"),
			likertQuestion("وقتی قولی می‌دهم، خودم را شخصاً موظف به انجام آن می‌دانم.", "responsibility"),
			likertQuestion("روزهای شلوغی که در آن‌ها پیشرفت مشخصی دارم به من انرژی می‌دهند.", "achiever"),
			likertQuestion("در شرایط نامطمئن، مسیرها و پیامدهای مختلف را با هم مقایسه می‌کنم.", "strategic"),
			likertQuestion("اغلب کتاب، دوره یا تجربه‌ای را انتخاب می‌کنم تا مهارت‌هایم رشد کنند.", "learner"),
			likertQuestion("افرادی که مرا خوب می‌شناسند، بیشترین انرژی و وفاداری من را دریافت می‌کنند.", "relator"),
			likertQuestion("می‌توانم اطلاعات ساده را جذاب و به‌یادماندنی بیان کنم.", "communication"),
			likertQuestion("دیگران روی من حساب می‌کنند چون وظایفم را جدی می‌گیرم.", "responsibility"),
		},
	}
}
