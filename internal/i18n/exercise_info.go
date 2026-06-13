package i18n

// info holds the extended exercise card text in [ru, en].
type info struct {
	Technique [2]string
	Mistake   [2]string
	Warning   [2]string // local safety warning (spec §20); empty = none
}

var exerciseInfo = map[string]info{
	"M01": {
		Technique: [2]string{"Ладони под плечами или чуть шире. Корпус одной линией, живот подтянут. Локти идут под 30–45° к корпусу. Не проваливай поясницу.", "Hands under or slightly wider than shoulders. Body in one line, belly tight. Elbows at 30–45°. Don't let the lower back sag."},
		Mistake:   [2]string{"Таз провисает, плечи уходят к ушам.", "Hips sag, shoulders creep toward the ears."},
	},
	"M02": {
		Technique: [2]string{"Ладони примерно под плечами, локти ближе к корпусу. Движение медленное. Не ставь ладони слишком близко — береги запястья.", "Hands about under the shoulders, elbows close to the body. Move slowly. Don't place hands too close — protect the wrists."},
		Mistake:   [2]string{"Локти резко уходят в стороны.", "Elbows flare out sharply."},
		Warning:   [2]string{"При боли в локтях или запястьях замени на обычные отжимания.", "If your elbows or wrists hurt, switch to regular push-ups."},
	},
	"M03": {
		Technique: [2]string{"Опустись вниз, задержка 1 секунда, поднимись без рывка. Корпус остаётся прямым.", "Lower down, hold 1 second, press up without a jerk. Keep the body straight."},
		Mistake:   [2]string{"Падение вниз без контроля.", "Dropping down without control."},
	},
	"C01": {
		Technique: [2]string{"Лёжа на спине, поясница прижата к полу. Ноги прямые или слегка согнуты. Скрещивающие движения умеренной амплитуды.", "Lie on your back, lower back pressed down. Legs straight or slightly bent. Cross the legs with a moderate range."},
		Mistake:   [2]string{"Поясница отрывается от пола.", "The lower back lifts off the floor."},
		Warning:   [2]string{"Если поясница отрывается — подними ноги выше или замени упражнение.", "If the lower back lifts — raise the legs higher or swap the exercise."},
	},
	"C02": {
		Technique: [2]string{"Локти под плечами, тело одной линией. Ягодицы слегка напряжены, живот подтянут, взгляд в пол.", "Elbows under shoulders, body in one line. Glutes lightly squeezed, belly tight, gaze down."},
		Mistake:   [2]string{"Провал поясницы или слишком высокий таз.", "Sagging lower back, or hips raised too high."},
	},
	"C03": {
		Technique: [2]string{"Лёжа на спине, поясница прижата. Ноги двигаются вверх-вниз короткой амплитудой, быстро, но контролируемо.", "Lie on your back, lower back pressed down. Legs move up and down in a short, quick but controlled range."},
		Mistake:   [2]string{"Ноги слишком низко, поясница выгибается.", "Legs too low, lower back arches."},
		Warning:   [2]string{"Не опускай ноги слишком низко.", "Don't lower the legs too far."},
	},
	"C04": {
		Technique: [2]string{"Лёжа на спине, руки вверх, ноги согнуты под 90°. Медленно опускай противоположную руку и ногу. Поясница прижата.", "Lie on your back, arms up, knees bent 90°. Slowly lower opposite arm and leg. Keep the lower back down."},
		Mistake:   [2]string{"Потеря контакта поясницы с полом.", "Losing lower-back contact with the floor."},
	},
	"C05": {
		Technique: [2]string{"Локоть под плечом, тело одной линией, таз не провисает. Свободная рука на поясе или вверх.", "Elbow under the shoulder, body in one line, hips up. Free hand on the waist or up."},
		Mistake:   [2]string{"Таз падает вниз.", "The hips drop."},
		Warning:   [2]string{"При дискомфорте в плече сократи время на 30–50%.", "If the shoulder is uncomfortable, cut the time by 30–50%."},
	},
	"C06": {
		Technique: [2]string{"Лёжа на спине, ноги в позиции 90/90. Поочерёдно касайся пяткой пола, поясница не отрывается.", "Lie on your back, legs 90/90. Tap the floor with each heel in turn; keep the lower back down."},
		Mistake:   [2]string{"Поясница отрывается.", "The lower back lifts."},
	},
	"C07": {
		Technique: [2]string{"Лёжа на спине, поясница прижата, ноги согнуты или слегка вытянуты, руки вдоль корпуса. Удерживай напряжение пресса.", "Lie on your back, lower back down, knees bent or slightly extended, arms by the body. Hold the abs tight."},
		Mistake:   [2]string{"Поясница отрывается.", "The lower back lifts."},
	},
	"C08": {
		Technique: [2]string{"Обычная планка на локтях. Сильно напряги ягодицы и пресс, локти как будто тяни к носкам. Держи коротко, но максимально качественно.", "A forearm plank. Squeeze glutes and abs hard, pull elbows toward the toes. Hold short but with top quality."},
		Mistake:   [2]string{"Пытаться держать как обычную планку 60–90 сек.", "Trying to hold it like a normal 60–90 s plank."},
		Warning:   [2]string{"Это короткая силовая планка — не держи её до отказа.", "This is a short power plank — don't hold it to failure."},
	},
	"B01": {
		Technique: [2]string{"Лёжа на животе, руки вперёд или вдоль корпуса. Слегка подними грудь и ноги, шея нейтрально, опускайся контролируемо.", "Lie face down, arms forward or by the body. Lift chest and legs slightly, neutral neck, lower with control."},
		Mistake:   [2]string{"Слишком сильный прогиб в пояснице.", "Over-arching the lower back."},
	},
	"B02": {
		Technique: [2]string{"Лёжа на животе, руки вдоль корпуса или в стороны. Движение как «ангелы на снегу» лицом вниз, активно работают лопатки.", "Lie face down, arms by the body or out. Move like 'snow angels' face down; the shoulder blades do the work."},
		Mistake:   [2]string{"Задирать голову вверх.", "Lifting the head up."},
	},
	"B03": {
		Technique: [2]string{"Лёжа на животе, слегка подними грудь и ноги, удерживай позицию. Дыхание не задерживай.", "Lie face down, lift chest and legs slightly, hold. Don't hold your breath."},
		Mistake:   [2]string{"Дискомфорт в пояснице — уменьши амплитуду.", "Lower-back discomfort — reduce the range."},
	},
	"B04": {
		Technique: [2]string{"Лёжа на животе, руки в форме буквы W. Подними локти и кисти от пола, сведи лопатки, опусти.", "Lie face down, arms in a 'W'. Lift elbows and hands, squeeze the shoulder blades, lower."},
		Mistake:   [2]string{"Движение идёт шеей, а не лопатками.", "Moving from the neck instead of the shoulder blades."},
	},
	"B05": {
		Technique: [2]string{"Лёжа на животе, руки в форме буквы Y. Подними руки невысоко, лопатки тяни вниз и назад, шея нейтрально.", "Lie face down, arms in a 'Y'. Lift the arms a little, pull the blades down and back, neutral neck."},
		Mistake:   [2]string{"Слишком большая амплитуда и прогиб поясницы.", "Too much range and lower-back arch."},
	},
	"G01": {
		Technique: [2]string{"Лёжа на спине, стопы на полу. Подними таз вверх, вверху напряги ягодицы, опусти таз контролируемо.", "Lie on your back, feet on the floor. Lift the hips, squeeze the glutes at the top, lower with control."},
		Mistake:   [2]string{"Переразгибание поясницы.", "Over-extending the lower back."},
	},
	"G02": {
		Technique: [2]string{"Подними таз, задержись 1–2 секунды сверху, опустись медленно.", "Lift the hips, hold 1–2 s at the top, lower slowly."},
		Mistake:   [2]string{"Переразгибание поясницы.", "Over-extending the lower back."},
	},
	"L01": {
		Technique: [2]string{"Стопы чуть шире таза, носки немного наружу. Таз назад, присед неглубокий (~45°), колени по линии носков, вес на середину стопы и пятку.", "Feet a bit wider than hips, toes slightly out. Hips back, shallow squat (~45°), knees track the toes, weight mid-foot and heel."},
		Mistake:   [2]string{"Колени заваливаются внутрь.", "Knees cave inward."},
		Warning:   [2]string{"Не уходи глубоко. При боли в колене замени на ягодичный мост.", "Don't go deep. If your knee hurts, switch to the glute bridge."},
	},
	"A01": {
		Technique: [2]string{"Правая рука сгибается в локте, левая давит на неё вниз, создавая сопротивление. Подъём медленный (3–4 сек). Затем поменяй стороны.", "The right arm curls; the left presses down on it for resistance. Raise slowly (3–4 s). Then switch sides."},
		Mistake:   [2]string{"Делать без сопротивления.", "Doing it without resistance."},
	},
	"A02": {
		Technique: [2]string{"Согни руку в локте на ~90°, второй рукой дави сверху вниз, рабочая рука сопротивляется. Удерживай напряжение.", "Bend the elbow to ~90°, press down with the other hand, the working arm resists. Hold the tension."},
		Mistake:   [2]string{"Задержка дыхания.", "Holding your breath."},
	},
	"S01": {
		Technique: [2]string{"Руки в стороны, маленькие круги вперёд, затем назад. Плечи не поднимай к ушам.", "Arms out to the sides, small circles forward then back. Don't raise the shoulders to the ears."},
		Mistake:   [2]string{"Слишком большая амплитуда и напряжение шеи.", "Too much range and neck tension."},
	},
}

func pick2(row [2]string, l Lang) string {
	if l == RU {
		return row[0]
	}
	return row[1]
}

// Technique / Mistake / Warning return localized extended-card text ("" if none).
func Technique(l Lang, id string) string {
	if i, ok := exerciseInfo[id]; ok {
		return pick2(i.Technique, l)
	}
	return ""
}
func Mistake(l Lang, id string) string {
	if i, ok := exerciseInfo[id]; ok {
		return pick2(i.Mistake, l)
	}
	return ""
}
func Warning(l Lang, id string) string {
	if i, ok := exerciseInfo[id]; ok {
		return pick2(i.Warning, l)
	}
	return ""
}
