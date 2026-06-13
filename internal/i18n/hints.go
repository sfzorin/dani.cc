package i18n

// hints holds the short on-screen technique cue per exercise ([ru, en]).
// Falls back to English; empty string means no hint shown.
var hints = map[string][2]string{
	"C02": {"Корпус одной линией, живот подтянут.", "Body in one line, belly tight."},
	"C01": {"Поясница прижата к полу.", "Lower back pressed to the floor."},
	"C03": {"Ноги двигаются коротко, поясница не отрывается.", "Short kicks, keep the lower back down."},
	"C04": {"Поясница прижата, движение медленное.", "Keep the lower back down, move slowly."},
	"C05": {"Таз не провисает, тело одной линией.", "Hips up, body in one line."},
	"C06": {"Поясница не отрывается от пола.", "Don't lift the lower back."},
	"C07": {"Поясница прижата к полу.", "Keep the lower back pressed down."},
	"C08": {"Короткая силовая планка — не держи до отказа.", "Short hard plank — don't hold to failure."},
	"M01": {"Локти под 30–45°, корпус ровный.", "Elbows 30–45°, body straight."},
	"M02": {"Локти ближе к корпусу, без боли в запястьях.", "Elbows close, no wrist pain."},
	"M03": {"Внизу пауза 1 секунда, без рывка.", "Pause 1 sec at the bottom, no bounce."},
	"B01": {"Подъём небольшой, шея расслаблена.", "Small lift, relax the neck."},
	"B02": {"Работают лопатки, голову не задирать.", "Move from the shoulder blades, don't lift the head."},
	"B03": {"Дыхание не задерживай.", "Keep breathing."},
	"G01": {"Толкайся пятками, напряги ягодицы сверху.", "Push through heels, squeeze glutes at the top."},
	"G02": {"Сверху пауза 1–2 секунды.", "Pause 1–2 sec at the top."},
	"L01": {"Неглубоко, колени по линии носков.", "Shallow, knees track the toes."},
	"A01": {"Одна рука сгибает, другая мешает.", "One arm curls, the other resists."},
	"A02": {"Держи угол 90°, не задерживай дыхание.", "Hold 90°, keep breathing."},
	"S01": {"Плечи вниз, круги маленькие.", "Shoulders down, small circles."},
	"W09": {"Спина прямая, наклон от таза.", "Straight back, hinge from the hips."},
}

// Hint returns the short technique cue for an exercise (empty if none).
func Hint(l Lang, id string) string {
	row, ok := hints[id]
	if !ok {
		return ""
	}
	if l == RU {
		return row[0]
	}
	return row[1]
}
