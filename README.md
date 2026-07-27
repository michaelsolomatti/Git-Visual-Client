🌳 Git Visual Client – Визуализация истории коммитов в терминале
Интерактивный консольный клиент Git с отображением графа коммитов, веток и тегов.
Реализован на 7 языках программирования – выберите свой!

✨ Возможности
📊 Визуализация графа коммитов – наглядное отображение структуры репозитория с использованием ASCII-символов.

🌿 Отображение веток и тегов – подсветка текущей ветки, показ всех ссылок.

🔍 Фильтрация и сортировка – опции --all для всех веток, --oneline для краткого вывода, --decorate для аннотаций.

⌨️ Интерактивный режим – навигация по коммитам с помощью клавиш (вверх/вниз), просмотр изменений (где реализовано).

🎨 Цветной вывод – выделение авторов, веток, коммитов.

⚡ Быстрая работа – эффективное использование памяти для больших репозиториев.

🛠 Интеграция с системой – использует нативный Git или библиотеки для работы с репозиториями.

📦 Поддерживаемые языки
Язык	Версия	Файл	Основная библиотека
Python	3.8+	git_visual.py	GitPython
Go	1.18+	git_visual.go	go-git
Rust	1.60+	git_visual.rs	git2
JavaScript	Node.js 14+	git_visual.js	simple-git
C#	.NET 6+	git_visual.cs	LibGit2Sharp
Java	11+	GitVisual.java	JGit
C++	C++17	git_visual.cpp	libgit2
🚀 Быстрый старт
1. Склонируйте репозиторий
bash
git clone https://github.com/yourname/git-visual-client.git
cd git-visual-client
2. Запустите на любом языке
Python

bash
pip install GitPython
python git_visual.py --graph --all --decorate
Go

bash
go mod init git_visual
go get github.com/go-git/go-git/v5
go run git_visual.go -graph -all -decorate
Rust (сборка)

bash
cargo new git_visual
# добавьте зависимости в Cargo.toml
cargo run -- --graph --all --decorate
JavaScript (Node.js)

bash
npm install simple-git
node git_visual.js --graph --all --decorate
C#

bash
dotnet new console -n git_visual
dotnet add package LibGit2Sharp
dotnet run -- --graph --all --decorate
Java (сборка с Maven/Gradle)

bash
javac -cp .:jgit.jar GitVisual.java
java -cp .:jgit.jar GitVisual --graph --all --decorate
C++ (сборка с libgit2)

bash
g++ -std=c++17 -I/usr/include/git2 git_visual.cpp -lgit2 -o git_visual
./git_visual --graph --all --decorate
📋 Пример вывода
Для репозитория с несколькими ветками программа выдаст что-то вроде:

text
* commit 3a4b5c6 (HEAD -> main) Добавлен новый функционал
| Author: John Doe <john@example.com>
| Date:   2025-03-15 10:30
|
* commit 1a2b3c4 (develop) Начало разработки
| Author: Jane Smith <jane@example.com>
| Date:   2025-03-14 09:00
|
| * commit 5e6f7g8 (feature/login) Добавлена авторизация
| | Author: John Doe <john@example.com>
| | Date:   2025-03-13 15:20
|/
* commit 9a0b1c2 Инициализация репозитория
  Author: Admin <admin@example.com>
  Date:   2025-03-12 12:00
⚙️ Опции командной строки
Флаг	Описание
--graph	Показать ASCII-граф коммитов
--all	Показать все ветки (включая удалённые)
--decorate	Показать ветки и теги рядом с коммитами
--oneline	Краткий вывод (только хеш и сообщение)
--author <имя>	Фильтр по автору
--since <дата>	Коммиты после указанной даты
--until <дата>	Коммиты до указанной даты
-n <число>	Ограничить количество коммитов
--color	Принудительно включить цвет
--help	Справка
Если запущено без опций, по умолчанию показывается граф текущей ветки с декором.

📄 Лицензия
MIT – свободно используйте, модифицируйте и распространяйте.

🤝 Вклад
Приветствуются pull request'ы! Если хотите добавить новый язык или улучшить существующий – создавайте issue.

🧠 Авторы
Проект создан в образовательных целях для демонстрации работы с Git на разных языках.

