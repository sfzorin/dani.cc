package i18n

// months[m-1] holds month names in Order = [ru, en, tr, de, es, fr, it].
var months = [12][7]string{
	{"Январь", "January", "Ocak", "Januar", "Enero", "Janvier", "Gennaio"},
	{"Февраль", "February", "Şubat", "Februar", "Febrero", "Février", "Febbraio"},
	{"Март", "March", "Mart", "März", "Marzo", "Mars", "Marzo"},
	{"Апрель", "April", "Nisan", "April", "Abril", "Avril", "Aprile"},
	{"Май", "May", "Mayıs", "Mai", "Mayo", "Mai", "Maggio"},
	{"Июнь", "June", "Haziran", "Juni", "Junio", "Juin", "Giugno"},
	{"Июль", "July", "Temmuz", "Juli", "Julio", "Juillet", "Luglio"},
	{"Август", "August", "Ağustos", "August", "Agosto", "Août", "Agosto"},
	{"Сентябрь", "September", "Eylül", "September", "Septiembre", "Septembre", "Settembre"},
	{"Октябрь", "October", "Ekim", "Oktober", "Octubre", "Octobre", "Ottobre"},
	{"Ноябрь", "November", "Kasım", "November", "Noviembre", "Novembre", "Novembre"},
	{"Декабрь", "December", "Aralık", "Dezember", "Diciembre", "Décembre", "Dicembre"},
}

// weekdays[i] holds short names, Monday-first (i=0 Monday .. 6 Sunday).
var weekdays = [7][7]string{
	{"Пн", "Mon", "Pzt", "Mo", "Lun", "Lun", "Lun"},
	{"Вт", "Tue", "Sal", "Di", "Mar", "Mar", "Mar"},
	{"Ср", "Wed", "Çar", "Mi", "Mié", "Mer", "Mer"},
	{"Чт", "Thu", "Per", "Do", "Jue", "Jeu", "Gio"},
	{"Пт", "Fri", "Cum", "Fr", "Vie", "Ven", "Ven"},
	{"Сб", "Sat", "Cmt", "Sa", "Sáb", "Sam", "Sab"},
	{"Вс", "Sun", "Paz", "So", "Dom", "Dim", "Dom"},
}

// Month returns the localized name of month m (1-12).
func Month(l Lang, m int) string {
	if m < 1 || m > 12 {
		return ""
	}
	i := index(l)
	if i < 0 {
		i = 1
	}
	return months[m-1][i]
}

// Weekday returns the localized short name for a Monday-first index (0-6).
func Weekday(l Lang, idx int) string {
	if idx < 0 || idx > 6 {
		return ""
	}
	i := index(l)
	if i < 0 {
		i = 1
	}
	return weekdays[idx][i]
}
