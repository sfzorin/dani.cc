package i18n

// exerciseNames maps exercise ID to names in Order = [ru, en, tr, de, es, fr, it].
var exerciseNames = map[string][7]string{
	// Warm-up
	"W01": {"Дыхание стоя", "Standing breathing", "Ayakta nefes", "Atmung im Stehen", "Respiración de pie", "Respiration debout", "Respirazione in piedi"},
	"W02": {"Повороты головы", "Head turns", "Baş çevirme", "Kopfdrehungen", "Giros de cabeza", "Rotations de la tête", "Rotazioni della testa"},
	"W03": {"Круги плечами назад", "Shoulder circles back", "Omuz çevirme (geri)", "Schulterkreisen rückwärts", "Círculos de hombro atrás", "Cercles d'épaules arrière", "Cerchi spalle indietro"},
	"W04": {"Круги плечами вперёд", "Shoulder circles forward", "Omuz çevirme (ileri)", "Schulterkreisen vorwärts", "Círculos de hombro adelante", "Cercles d'épaules avant", "Cerchi spalle avanti"},
	"W05": {"Круги руками вперёд", "Arm circles forward", "Kol çevirme (ileri)", "Armkreisen vorwärts", "Círculos de brazos adelante", "Cercles des bras avant", "Cerchi braccia avanti"},
	"W06": {"Круги руками назад", "Arm circles back", "Kol çevirme (geri)", "Armkreisen rückwärts", "Círculos de brazos atrás", "Cercles des bras arrière", "Cerchi braccia indietro"},
	"W07": {"Сведение лопаток", "Scapula squeeze", "Kürek kemiği sıkma", "Schulterblätter zusammen", "Juntar omóplatos", "Serrage des omoplates", "Avvicinare le scapole"},
	"W08": {"Повороты корпуса", "Torso twists", "Gövde dönüşü", "Rumpfdrehung", "Giros de torso", "Rotation du buste", "Torsioni del busto"},
	"W09": {"Наклон с прямой спиной", "Hip hinge", "Kalça menteşesi", "Hüftbeuge", "Bisagra de cadera", "Charnière de hanche", "Hip hinge"},
	"W10": {"Подъём на носки", "Calf raises", "Topuk kaldırma", "Wadenheben", "Elevación de pantorrillas", "Mollets debout", "Sollevamento polpacci"},

	// Main
	"M01": {"Отжимания", "Push-ups", "Şınav", "Liegestütze", "Flexiones", "Pompes", "Piegamenti"},
	"M02": {"Узкие отжимания", "Narrow push-ups", "Dar şınav", "Enge Liegestütze", "Flexiones cerradas", "Pompes serrées", "Piegamenti stretti"},
	"M03": {"Отжимания с паузой", "Paused push-ups", "Duraklamalı şınav", "Liegestütze mit Pause", "Flexiones con pausa", "Pompes avec pause", "Piegamenti con pausa"},
	"C01": {"Ножницы", "Scissors", "Makas", "Schere", "Tijeras", "Ciseaux", "Forbici"},
	"C02": {"Планка", "Plank", "Plank", "Unterarmstütz", "Plancha", "Planche", "Plank"},
	"C03": {"Флаттер", "Flutter kicks", "Çırpma hareketi", "Flutter Kicks", "Patada de aleteo", "Battements de jambes", "Flutter kicks"},
	"C04": {"Дэд-баг", "Dead bug", "Ölü böcek", "Dead Bug", "Bicho muerto", "Dead bug", "Dead bug"},
	"C05": {"Боковая планка", "Side plank", "Yan plank", "Seitstütz", "Plancha lateral", "Planche latérale", "Plank laterale"},
	"C06": {"Toe taps", "Toe taps", "Ayak ucu dokunuşu", "Toe Taps", "Toques de pie", "Touches d'orteils", "Toe taps"},
	"C07": {"Hollow hold", "Hollow hold", "Hollow hold", "Hollow Hold", "Hollow hold", "Hollow hold", "Hollow hold"},
	"C08": {"RKC-планка", "RKC plank", "RKC plank", "RKC-Plank", "Plancha RKC", "Planche RKC", "Plank RKC"},
	"B01": {"Лодочка", "Boat", "Kayık", "Boot", "Barquito", "Bateau", "Barchetta"},
	"B02": {"Обратные снежные ангелы", "Reverse snow angels", "Ters kar meleği", "Reverse Snow Angels", "Ángeles invertidos", "Anges inversés", "Angeli inversi"},
	"B03": {"Лодочка (удержание)", "Boat hold", "Kayık tutuş", "Boot halten", "Barquito isométrico", "Bateau tenu", "Barchetta tenuta"},
	"B04": {"W-подъём лёжа", "Prone W-raise", "W kaldırma", "W-Heben", "Elevación en W", "Levée en W", "Sollevamento a W"},
	"B05": {"Y-подъём лёжа", "Prone Y-raise", "Y kaldırma", "Y-Heben", "Elevación en Y", "Levée en Y", "Sollevamento a Y"},
	"G01": {"Ягодичный мост", "Glute bridge", "Köprü", "Glute Bridge", "Puente de glúteos", "Pont fessier", "Ponte glutei"},
	"G02": {"Мост с паузой", "Glute bridge hold", "Duraklamalı köprü", "Glute Bridge mit Pause", "Puente con pausa", "Pont avec pause", "Ponte con pausa"},
	"L01": {"Мини-присед", "Mini squat", "Mini çömelme", "Mini-Kniebeuge", "Sentadilla mini", "Mini squat", "Mini squat"},
	"A01": {"Бицепс с самосопротивлением", "Bicep self-resistance", "Pazı (kendine direnç)", "Bizeps mit Eigenwiderstand", "Bíceps autorresistencia", "Biceps auto-résistance", "Bicipiti auto-resistenza"},
	"A02": {"Изометрия бицепса", "Bicep isometric hold", "Pazı izometrik", "Bizeps isometrisch", "Bíceps isométrico", "Biceps isométrique", "Bicipiti isometrico"},
	"S01": {"Круги руками", "Arm circles", "Kol çevirme", "Armkreisen", "Círculos de brazos", "Cercles des bras", "Cerchi con le braccia"},

	// Cool-down
	"CD01": {"Дыхание лёжа", "Lying breathing", "Sırtüstü nefes", "Atmung im Liegen", "Respiración tumbado", "Respiration allongée", "Respirazione sdraiato"},
	"CD02": {"Растяжка груди", "Chest stretch", "Göğüs esnetme", "Brustdehnung", "Estiramiento de pecho", "Étirement pectoral", "Allungamento pettorali"},
	"CD03": {"Скручивание лёжа", "Lying twist", "Sırtüstü dönüş", "Liegende Drehung", "Torsión tumbado", "Torsion allongée", "Torsione sdraiato"},
	"CD04": {"Растяжка бедра", "Hamstring stretch", "Arka bacak esnetme", "Beinrückseite dehnen", "Estiramiento de isquios", "Étirement ischios", "Allungamento femorali"},
	"CD05": {"Растяжка «фигура 4»", "Figure-4 glute stretch", "Figür-4 esneme", "Figure-4-Dehnung", "Estiramiento figura 4", "Étirement figure 4", "Stretching figura 4"},
	"CD06": {"Поза сфинкса", "Sphinx pose", "Sfenks pozu", "Sphinx-Haltung", "Postura de esfinge", "Posture du sphinx", "Posizione della sfinge"},
}
