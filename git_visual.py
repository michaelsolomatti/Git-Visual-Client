# git_visual.py
import sys
import argparse
from git import Repo, Commit, Actor
from datetime import datetime
import os

COLORS = {
    'red': '\033[91m',
    'green': '\033[92m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'cyan': '\033[96m',
    'reset': '\033[0m',
}

class GitVisual:
    def __init__(self, repo_path='.', color=True, graph=True, all_branches=False,
                 decorate=True, oneline=False, author=None, since=None, until=None, max_count=None):
        self.repo_path = repo_path
        self.color = color
        self.graph = graph
        self.all = all_branches
        self.decorate = decorate
        self.oneline = oneline
        self.author = author
        self.since = since
        self.until = until
        self.max_count = max_count
        self.repo = None

    def open_repo(self):
        try:
            self.repo = Repo(self.repo_path)
        except Exception as e:
            print(f"Error opening repository: {e}", file=sys.stderr)
            sys.exit(1)

    def colorize(self, text, color):
        if self.color:
            return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"
        return text

    def get_commits(self):
        # Строим опции для rev-list
        args = []
        if self.all:
            args.append('--all')
        if self.author:
            args.append(f'--author={self.author}')
        if self.since:
            args.append(f'--since={self.since}')
        if self.until:
            args.append(f'--until={self.until}')
        if self.max_count:
            args.append(f'-n {self.max_count}')
        # Получаем коммиты
        revs = self.repo.git.rev_list('--pretty=format:%h|%an|%ae|%ad|%s', *args).splitlines()
        commits = []
        for line in revs:
            if line.startswith('commit '):
                continue
            parts = line.split('|', 4)
            if len(parts) >= 5:
                h, author, email, date, msg = parts[0], parts[1], parts[2], parts[3], parts[4]
                commits.append({
                    'hash': h,
                    'author': author,
                    'email': email,
                    'date': date,
                    'message': msg
                })
        return commits

    def print_commit(self, commit, prefix='', last=True, branches=None):
        if branches is None:
            branches = {}
        connector = '└── ' if last else '├── '
        hash_color = 'yellow' if self.color else ''
        author_color = 'cyan' if self.color else ''
        date_color = 'green' if self.color else ''
        msg_color = 'reset'

        # Определяем ветки для этого коммита
        refs = branches.get(commit['hash'], [])
        ref_str = ' (' + ', '.join(refs) + ')' if refs and self.decorate else ''

        if self.oneline:
            line = f"{connector}{self.colorize(commit['hash'], hash_color)} {self.colorize(commit['message'], msg_color)}{self.colorize(ref_str, 'blue')}"
            print(prefix + line)
        else:
            print(f"{prefix}{connector}{self.colorize(commit['hash'], hash_color)}{self.colorize(ref_str, 'blue')}")
            print(f"{prefix}    {self.colorize('Author:', author_color)} {self.colorize(commit['author'], author_color)} <{commit['email']}>")
            print(f"{prefix}    {self.colorize('Date:', date_color)} {self.colorize(commit['date'], date_color)}")
            print(f"{prefix}    {self.colorize(commit['message'], msg_color)}")
            print()

    def show_graph(self):
        # Получаем коммиты
        commits = self.get_commits()
        if not commits:
            print("No commits found.")
            return
        # Собираем информацию о ветках для декора
        branches = {}
        if self.decorate:
            for ref in self.repo.references:
                if ref.is_detached:
                    continue
                try:
                    commit_hash = ref.commit.hexsha[:7]
                    branches.setdefault(commit_hash, []).append(ref.name)
                except:
                    pass
            # Добавляем HEAD
            if self.repo.head.is_detached:
                try:
                    h = self.repo.head.commit.hexsha[:7]
                    branches.setdefault(h, []).append('HEAD')
                except:
                    pass
            else:
                try:
                    h = self.repo.head.commit.hexsha[:7]
                    branches.setdefault(h, []).append('HEAD -> ' + self.repo.active_branch.name)
                except:
                    pass

        # Печатаем каждый коммит с отступами
        for i, c in enumerate(commits):
            last = (i == len(commits) - 1)
            self.print_commit(c, prefix='', last=last, branches=branches)

def main():
    parser = argparse.ArgumentParser(description="Git Visual Client")
    parser.add_argument('--graph', action='store_true', default=True, help='Show ASCII graph')
    parser.add_argument('--all', action='store_true', help='Show all branches')
    parser.add_argument('--decorate', action='store_true', default=True, help='Show refs')
    parser.add_argument('--oneline', action='store_true', help='Compact output')
    parser.add_argument('--author', help='Filter by author')
    parser.add_argument('--since', help='Commits after date')
    parser.add_argument('--until', help='Commits before date')
    parser.add_argument('-n', '--max-count', type=int, help='Limit number of commits')
    parser.add_argument('--color', action='store_true', default=True, help='Color output')
    parser.add_argument('--path', default='.', help='Repository path')
    args = parser.parse_args()

    visual = GitVisual(
        repo_path=args.path,
        color=args.color,
        graph=args.graph,
        all_branches=args.all,
        decorate=args.decorate,
        oneline=args.oneline,
        author=args.author,
        since=args.since,
        until=args.until,
        max_count=args.max_count
    )
    visual.open_repo()
    visual.show_graph()

if __name__ == '__main__':
    main()
