package i18n

// ui maps a string key to translations in Order = [ru, en, tr, de, es, fr, it].
var ui = map[string][7]string{
	// App / brand
	"app.name":    {"Зарядка", "Morning", "Sabah", "Morgen", "Mañana", "Matin", "Mattina"},
	"app.tagline": {"Утренняя зарядка за 30 дней", "30-day morning workout", "30 günlük sabah egzersizi", "30-Tage-Morgenworkout", "Rutina matutina de 30 días", "Routine matinale de 30 jours", "Allenamento mattutino di 30 giorni"},

	// Guest mode
	"guest.note": {"Вы занимаетесь как гость. Зарегистрируйтесь, чтобы не потерять прогресс.", "You're a guest. Sign up so you don't lose your progress.", "Misafir olarak çalışıyorsunuz. İlerlemeyi kaybetmemek için kayıt olun.", "Du trainierst als Gast. Registriere dich, um deinen Fortschritt zu sichern.", "Entrenas como invitado. Regístrate para no perder tu progreso.", "Vous êtes invité. Inscrivez-vous pour sauvegarder votre progression.", "Ti alleni come ospite. Registrati per non perdere i progressi."},

	// Nav
	"nav.home":     {"Главная", "Home", "Ana sayfa", "Start", "Inicio", "Accueil", "Home"},
	"nav.settings": {"Настройки", "Settings", "Ayarlar", "Einstellungen", "Ajustes", "Réglages", "Impostazioni"},
	"nav.logout":   {"Выйти", "Log out", "Çıkış", "Abmelden", "Salir", "Déconnexion", "Esci"},

	// Auth
	"auth.login":          {"Вход", "Log in", "Giriş", "Anmelden", "Iniciar sesión", "Connexion", "Accedi"},
	"auth.register":       {"Регистрация", "Sign up", "Kayıt ol", "Registrieren", "Registrarse", "Inscription", "Registrati"},
	"auth.name":           {"Имя", "Name", "İsim", "Name", "Nombre", "Prénom", "Nome"},
	"auth.email":          {"Эл. почта", "Email", "E-posta", "E-Mail", "Correo", "E-mail", "Email"},
	"auth.password":       {"Пароль", "Password", "Parola", "Passwort", "Contraseña", "Mot de passe", "Password"},
	"auth.login_btn":      {"Войти", "Log in", "Giriş yap", "Anmelden", "Entrar", "Se connecter", "Accedi"},
	"auth.register_btn":   {"Создать аккаунт", "Create account", "Hesap oluştur", "Konto erstellen", "Crear cuenta", "Créer un compte", "Crea account"},
	"auth.have_account":   {"Уже есть аккаунт?", "Already have an account?", "Hesabın var mı?", "Schon ein Konto?", "¿Ya tienes cuenta?", "Déjà un compte ?", "Hai già un account?"},
	"auth.continue_guest": {"Продолжить без регистрации", "Continue without signing up", "Kayıt olmadan devam et", "Ohne Anmeldung fortfahren", "Continuar sin registrarse", "Continuer sans inscription", "Continua senza registrarti"},
	"auth.no_account":     {"Нет аккаунта?", "No account yet?", "Hesabın yok mu?", "Noch kein Konto?", "¿Sin cuenta?", "Pas de compte ?", "Non hai un account?"},
	"auth.err_required":   {"Заполните все поля", "Fill in all fields", "Tüm alanları doldurun", "Bitte alle Felder ausfüllen", "Completa todos los campos", "Remplissez tous les champs", "Compila tutti i campi"},
	"auth.err_invalid":    {"Неверная почта или пароль", "Wrong email or password", "Hatalı e-posta veya parola", "Falsche E-Mail oder Passwort", "Correo o contraseña incorrectos", "E-mail ou mot de passe incorrect", "Email o password errati"},
	"auth.err_exists":     {"Такой пользователь уже есть", "This account already exists", "Bu hesap zaten var", "Konto existiert bereits", "La cuenta ya existe", "Ce compte existe déjà", "Account già esistente"},
	"auth.err_short_pass": {"Пароль не короче 6 символов", "Password must be 6+ characters", "Parola en az 6 karakter", "Passwort mind. 6 Zeichen", "Mínimo 6 caracteres", "6 caractères minimum", "Almeno 6 caratteri"},
	"auth.err_bad_email":  {"Введите корректную почту", "Enter a valid email", "Geçerli bir e-posta girin", "Gültige E-Mail eingeben", "Introduce un correo válido", "Entrez un e-mail valide", "Inserisci un'email valida"},

	// Home
	"home.morning":       {"Доброе утро", "Good morning", "Günaydın", "Guten Morgen", "Buenos días", "Bonjour", "Buongiorno"},
	"home.day":          {"День", "Day", "Gün", "Tag", "Día", "Jour", "Giorno"},
	"home.of":           {"из", "of", "/", "von", "de", "sur", "di"},
	"home.start":        {"Начать", "Start", "Başla", "Starten", "Empezar", "Commencer", "Inizia"},
	"home.continue":     {"Продолжить", "Continue", "Devam et", "Weiter", "Continuar", "Continuer", "Continua"},
	"home.repeat":       {"Повторить", "Repeat", "Tekrarla", "Wiederholen", "Repetir", "Refaire", "Ripeti"},
	"home.streak":       {"Серия", "Streak", "Seri", "Serie", "Racha", "Série", "Serie"},
	"home.days":         {"дн.", "days", "gün", "Tage", "días", "jours", "giorni"},
	"home.done_today":   {"Сегодня выполнено! 🎉", "Done for today! 🎉", "Bugün tamam! 🎉", "Heute geschafft! 🎉", "¡Hecho por hoy! 🎉", "Terminé pour aujourd'hui ! 🎉", "Fatto per oggi! 🎉"},
	"home.calendar":     {"Календарь", "Calendar", "Takvim", "Kalender", "Calendario", "Calendrier", "Calendario"},
	"home.est":          {"≈ время", "≈ time", "≈ süre", "≈ Zeit", "≈ tiempo", "≈ durée", "≈ tempo"},
	"home.min":          {"мин", "min", "dk", "Min", "min", "min", "min"},
	"home.exercises":    {"упражнений", "exercises", "egzersiz", "Übungen", "ejercicios", "exercices", "esercizi"},
	"home.cycle":        {"Цикл", "Cycle", "Döngü", "Zyklus", "Ciclo", "Cycle", "Ciclo"},
	"home.diff":         {"Сложность растёт", "Difficulty increases", "Zorluk artar", "Schwierigkeit steigt", "La dificultad aumenta", "La difficulté augmente", "La difficoltà aumenta"},
	"home.this_month":   {"Этот месяц", "This month", "Bu ay", "Dieser Monat", "Este mes", "Ce mois", "Questo mese"},

	// Workout player
	"workout.get_ready": {"Приготовься", "Get ready", "Hazır ol", "Mach dich bereit", "Prepárate", "Préparez-vous", "Preparati"},
	"workout.round":     {"Круг", "Round", "Tur", "Runde", "Ronda", "Tour", "Round"},
	"workout.rest":      {"Отдых", "Rest", "Dinlen", "Pause", "Descanso", "Repos", "Riposo"},
	"workout.next":      {"Далее", "Next", "Sıradaki", "Weiter", "Siguiente", "Suivant", "Prossimo"},
	"workout.pause":     {"Пауза", "Pause", "Duraklat", "Pause", "Pausa", "Pause", "Pausa"},
	"workout.resume":    {"Продолжить", "Resume", "Devam", "Fortsetzen", "Reanudar", "Reprendre", "Riprendi"},
	"workout.skip":      {"Пропустить", "Skip", "Atla", "Überspringen", "Saltar", "Passer", "Salta"},
	"workout.prev":      {"Назад", "Back", "Geri", "Zurück", "Atrás", "Précédent", "Indietro"},
	"workout.done":      {"Готово", "Done", "Tamam", "Fertig", "Listo", "Fait", "Fatto"},
	"workout.warmup":    {"Разминка", "Warm-up", "Isınma", "Aufwärmen", "Calentamiento", "Échauffement", "Riscaldamento"},
	"workout.main":      {"Основная часть", "Main set", "Ana bölüm", "Hauptteil", "Parte principal", "Partie principale", "Parte principale"},
	"workout.cooldown":  {"Заминка", "Cool-down", "Soğuma", "Abwärmen", "Enfriamiento", "Retour au calme", "Defaticamento"},
	"workout.quit":      {"Выйти", "Quit", "Çık", "Beenden", "Salir", "Quitter", "Esci"},
	"workout.x_reps":    {"повторений", "reps", "tekrar", "Wdh.", "reps", "reps", "ripet."},
	"workout.sec":       {"сек", "sec", "sn", "Sek", "seg", "s", "sec"},
	"workout.breaths":   {"дыханий", "breaths", "nefes", "Atemzüge", "respiraciones", "respirations", "respiri"},
	"workout.side":      {"сторона", "side", "taraf", "Seite", "lado", "côté", "lato"},

	// Exercise info card
	"info.technique": {"Техника", "Technique", "Teknik", "Technik", "Técnica", "Technique", "Tecnica"},
	"info.mistake":   {"Частая ошибка", "Common mistake", "Sık hata", "Häufiger Fehler", "Error común", "Erreur fréquente", "Errore comune"},
	"info.warning":   {"Внимание", "Caution", "Dikkat", "Achtung", "Atención", "Attention", "Attenzione"},
	"info.close":     {"Закрыть", "Close", "Kapat", "Schließen", "Cerrar", "Fermer", "Chiudi"},

	// Safety disclaimer (spec §20)
	"safety.title": {"Безопасность", "Safety", "Güvenlik", "Sicherheit", "Seguridad", "Sécurité", "Sicurezza"},
	"safety.text":  {"Программа не заменяет консультацию врача. При боли в груди, головокружении, острой боли в колене, пояснице, плече или локте остановите тренировку. Если боль повторяется — снизьте нагрузку или обратитесь к специалисту.", "This program is not a substitute for medical advice. Stop if you feel chest pain, dizziness, or sharp pain in a knee, lower back, shoulder or elbow. If pain recurs, reduce the load or see a specialist.", "Bu program tıbbi tavsiyenin yerine geçmez. Göğüs ağrısı, baş dönmesi veya dizde, belde, omuzda ya da dirsekte keskin ağrı olursa durun. Ağrı tekrarlarsa yükü azaltın veya uzmana başvurun.", "Dieses Programm ersetzt keine ärztliche Beratung. Bei Brustschmerzen, Schwindel oder stechenden Schmerzen in Knie, Rücken, Schulter oder Ellenbogen abbrechen. Bei wiederkehrenden Schmerzen die Belastung senken oder einen Facharzt aufsuchen.", "Este programa no sustituye el consejo médico. Detente si sientes dolor en el pecho, mareo o dolor agudo en rodilla, espalda, hombro o codo. Si el dolor se repite, reduce la carga o consulta a un especialista.", "Ce programme ne remplace pas un avis médical. Arrêtez en cas de douleur thoracique, de vertige ou de douleur vive au genou, au dos, à l'épaule ou au coude. Si la douleur persiste, réduisez la charge ou consultez un spécialiste.", "Questo programma non sostituisce il parere medico. Fermati in caso di dolore al petto, vertigini o dolore acuto a ginocchio, schiena, spalla o gomito. Se il dolore si ripete, riduci il carico o consulta uno specialista."},
	"safety.ok":    {"Понятно, начать", "Got it, start", "Anladım, başla", "Verstanden, los", "Entendido, empezar", "Compris, commencer", "Capito, inizia"},

	// Post-workout feedback
	"fb.title":      {"Как прошло?", "How did it go?", "Nasıl geçti?", "Wie war es?", "¿Cómo fue?", "Comment c'était ?", "Com'è andata?"},
	"fb.difficulty": {"Сложность (1–10)", "Difficulty (1–10)", "Zorluk (1–10)", "Schwierigkeit (1–10)", "Dificultad (1–10)", "Difficulté (1–10)", "Difficoltà (1–10)"},
	"fb.knee":       {"Колено", "Knee", "Diz", "Knie", "Rodilla", "Genou", "Ginocchio"},
	"fb.back":       {"Поясница", "Lower back", "Bel", "Unterer Rücken", "Zona lumbar", "Bas du dos", "Zona lombare"},
	"fb.shoulder":   {"Плечи / локти / запястья", "Shoulders / elbows / wrists", "Omuz / dirsek / bilek", "Schultern / Ellbogen / Handgelenke", "Hombros / codos / muñecas", "Épaules / coudes / poignets", "Spalle / gomiti / polsi"},
	"fb.energy":     {"Энергия после", "Energy after", "Sonrası enerji", "Energie danach", "Energía después", "Énergie après", "Energia dopo"},
	"fb.none":       {"Нет", "None", "Yok", "Keine", "Nada", "Aucune", "Nessuno"},
	"fb.some":       {"Немного", "A bit", "Biraz", "Etwas", "Un poco", "Un peu", "Un po'"},
	"fb.pain":       {"Больно", "Pain", "Ağrı", "Schmerz", "Dolor", "Douleur", "Dolore"},
	"fb.worse":      {"Хуже", "Worse", "Daha kötü", "Schlechter", "Peor", "Pire", "Peggio"},
	"fb.normal":     {"Норм", "OK", "Normal", "OK", "Normal", "OK", "OK"},
	"fb.better":     {"Лучше", "Better", "Daha iyi", "Besser", "Mejor", "Mieux", "Meglio"},
	"fb.save":       {"Сохранить", "Save", "Kaydet", "Speichern", "Guardar", "Enregistrer", "Salva"},
	"fb.skip":       {"Пропустить", "Skip", "Atla", "Überspringen", "Saltar", "Passer", "Salta"},

	// Completion
	"done.title":     {"Отлично! 👏", "Awesome! 👏", "Harika! 👏", "Super! 👏", "¡Genial! 👏", "Bravo ! 👏", "Fantastico! 👏"},
	"done.bravo":     {"ты молодец! 👏", "you did great! 👏", "harikasın! 👏", "stark gemacht! 👏", "¡lo hiciste genial! 👏", "bravo à toi ! 👏", "sei grande! 👏"},
	"done.streak":    {"Серия", "Streak", "Seri", "Serie", "Racha", "Série", "Serie"},
	"done.back_home": {"На главную", "Back home", "Ana sayfaya", "Zur Startseite", "Al inicio", "Accueil", "Alla home"},
	"done.day_done":  {"День завершён", "Day complete", "Gün tamamlandı", "Tag geschafft", "Día completado", "Jour terminé", "Giorno completato"},

	// Settings
	"settings.title":     {"Настройки", "Settings", "Ayarlar", "Einstellungen", "Ajustes", "Réglages", "Impostazioni"},
	"settings.rest":      {"Пауза между упражнениями", "Rest between exercises", "Egzersizler arası dinlenme", "Pause zwischen Übungen", "Descanso entre ejercicios", "Repos entre exercices", "Riposo tra esercizi"},
	"settings.language":  {"Язык", "Language", "Dil", "Sprache", "Idioma", "Langue", "Lingua"},
	"settings.voice":     {"Голосовые подсказки", "Voice cues", "Sesli komutlar", "Sprachansagen", "Indicaciones de voz", "Indications vocales", "Indicazioni vocali"},
	"settings.on":        {"Вкл", "On", "Açık", "An", "Sí", "Oui", "Sì"},
	"settings.off":       {"Выкл", "Off", "Kapalı", "Aus", "No", "Non", "No"},
	"settings.save":      {"Сохранить", "Save", "Kaydet", "Speichern", "Guardar", "Enregistrer", "Salva"},
	"settings.saved":     {"Сохранено", "Saved", "Kaydedildi", "Gespeichert", "Guardado", "Enregistré", "Salvato"},
	"settings.account":   {"Аккаунт", "Account", "Hesap", "Konto", "Cuenta", "Compte", "Account"},
	"settings.reset":     {"Сбросить прогресс", "Reset progress", "İlerlemeyi sıfırla", "Fortschritt zurücksetzen", "Reiniciar progreso", "Réinitialiser", "Azzera progressi"},

	// Voice cues (spoken)
	"cue.ready":       {"Приготовься", "Get ready", "Hazır ol", "Mach dich bereit", "Prepárate", "Préparez-vous", "Preparati"},
	"cue.go":          {"Начали", "Go", "Başla", "Los", "Vamos", "C'est parti", "Via"},
	"cue.rest":        {"Отдых", "Rest", "Dinlen", "Pause", "Descanso", "Repos", "Riposo"},
	"cue.next":        {"Дальше", "Next up", "Sıradaki", "Als Nächstes", "A continuación", "Ensuite", "Prossimo"},
	"cue.halfway":     {"Половина пройдена", "Halfway there", "Yarısı bitti", "Halbzeit", "A mitad", "À mi-chemin", "A metà"},
	"cue.last":        {"Последнее упражнение", "Last exercise", "Son egzersiz", "Letzte Übung", "Último ejercicio", "Dernier exercice", "Ultimo esercizio"},
	"cue.well_done":   {"Молодец", "Well done", "Aferin", "Gut gemacht", "Bien hecho", "Bien joué", "Bravo"},
	"cue.switch_side": {"Другая сторона", "Other side", "Diğer taraf", "Andere Seite", "Otro lado", "Autre côté", "Altro lato"},
	"cue.workout_end": {"Тренировка завершена", "Workout complete", "Antrenman bitti", "Workout fertig", "Entrenamiento completado", "Séance terminée", "Allenamento finito"},

	// Random encouragement shown/spoken at the end (one picked at random)
	"enc.0": {"Ты сделал это!", "You did it!", "Başardın!", "Du hast es geschafft!", "¡Lo lograste!", "Vous l'avez fait !", "Ce l'hai fatta!"},
	"enc.1": {"Сильно! Так держать!", "Strong! Keep it up!", "Güçlü! Devam et!", "Stark! Weiter so!", "¡Fuerte! ¡Sigue así!", "Fort ! Continuez !", "Forte! Continua così!"},
	"enc.2": {"Ещё один день в копилку!", "Another day in the bank!", "Bir gün daha tamam!", "Wieder ein Tag geschafft!", "¡Otro día más!", "Encore un jour de gagné !", "Un altro giorno fatto!"},
	"enc.3": {"Тело скажет спасибо!", "Your body thanks you!", "Vücudun teşekkür ediyor!", "Dein Körper dankt dir!", "¡Tu cuerpo te lo agradece!", "Votre corps vous remercie !", "Il tuo corpo ti ringrazia!"},
	"enc.4": {"Дисциплина — это свобода!", "Discipline is freedom!", "Disiplin özgürlüktür!", "Disziplin ist Freiheit!", "¡La disciplina es libertad!", "La discipline, c'est la liberté !", "La disciplina è libertà!"},
	"enc.5": {"Завтра будет ещё легче!", "Tomorrow will be even easier!", "Yarın daha kolay olacak!", "Morgen wird's noch leichter!", "¡Mañana será más fácil!", "Demain sera plus facile !", "Domani sarà ancora più facile!"},
}
