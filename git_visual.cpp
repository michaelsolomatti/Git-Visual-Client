// git_visual.cpp
#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <cstring>
#include <git2.h>
#include <unistd.h>

using namespace std;

bool colorEnabled = true;

string colorize(const string& text, const string& color) {
    if (!colorEnabled) return text;
    map<string, string> codes = {
        {"red", "\033[91m"},
        {"green", "\033[92m"},
        {"yellow", "\033[93m"},
        {"blue", "\033[94m"},
        {"cyan", "\033[96m"},
        {"reset", "\033[0m"}
    };
    return codes[color] + text + codes["reset"];
}

void printCommit(git_commit* commit, const string& prefix, bool last, map<string, vector<string>>& refs, bool decorate, bool oneline) {
    const git_oid* oid = git_commit_id(commit);
    char hash[8];
    git_oid_tostr(hash, sizeof(hash), oid);
    string hashStr(hash, 7);

    string refStr;
    if (decorate && refs.count(hashStr)) {
        refStr = " (";
        for (size_t i = 0; i < refs[hashStr].size(); ++i) {
            if (i) refStr += ", ";
            refStr += refs[hashStr][i];
        }
        refStr += ")";
    }

    string connector = last ? "└── " : "├── ";
    const char* msg = git_commit_message(commit);
    string message = msg ? msg : "";
    size_t newline = message.find('\n');
    if (newline != string::npos) message = message.substr(0, newline);

    if (oneline) {
        cout << prefix << connector << colorize(hashStr, "yellow") << colorize(refStr, "blue") << " " << colorize(message, "reset") << endl;
    } else {
        cout << prefix << connector << colorize(hashStr, "yellow") << colorize(refStr, "blue") << endl;
        const git_signature* author = git_commit_author(commit);
        string date = ctime(&author->when.time);
        date.pop_back();
        cout << prefix << "    " << colorize("Author:", "cyan") << " " << author->name << " <" << author->email << ">" << endl;
        cout << prefix << "    " << colorize("Date:", "green") << " " << date << endl;
        cout << prefix << "    " << colorize(message, "reset") << "\n" << endl;
    }
}

int main(int argc, char* argv[]) {
    git_libgit2_init();
    string repoPath = ".";
    bool graph = true, all = false, decorate = true, oneline = false;
    string authorFilter;
    int maxCount = 0;
    colorEnabled = isatty(STDOUT_FILENO); // автоопределение

    for (int i = 1; i < argc; ++i) {
        string arg = argv[i];
        if (arg == "--path" && i+1 < argc) repoPath = argv[++i];
        else if (arg == "--graph" && i+1 < argc) graph = (string(argv[++i]) == "true");
        else if (arg == "--all") all = true;
        else if (arg == "--decorate" && i+1 < argc) decorate = (string(argv[++i]) == "true");
        else if (arg == "--oneline") oneline = true;
        else if (arg == "--author" && i+1 < argc) authorFilter = argv[++i];
        else if (arg == "-n" && i+1 < argc) maxCount = stoi(argv[++i]);
        else if (arg == "--color" && i+1 < argc) colorEnabled = (string(argv[++i]) == "true");
        else if (arg == "--help") {
            cout << "Usage: ..." << endl;
            return 0;
        }
    }

    git_repository* repo = nullptr;
    int err = git_repository_open(&repo, repoPath.c_str());
    if (err < 0) {
        cerr << "Failed to open repo: " << git_error_last()->message << endl;
        return 1;
    }

    // Получаем список коммитов
    git_revwalk* walk;
    git_revwalk_new(&walk, repo);
    if (all) {
        git_revwalk_push_glob(walk, "refs/*");
    } else {
        git_oid head_oid;
        git_reference* head;
        git_repository_head(&head, repo);
        git_reference_target(&head_oid, head);
        git_revwalk_push(walk, &head_oid);
    }
    git_revwalk_sorting(walk, GIT_SORT_TIME);

    // Собираем ссылки для декора
    map<string, vector<string>> refs;
    if (decorate) {
        git_strarray ref_list;
        git_reference_list(&ref_list, repo);
        for (size_t i = 0; i < ref_list.count; ++i) {
            git_reference* ref;
            if (git_reference_lookup(&ref, repo, ref_list.strings[i]) == 0) {
                const git_oid* oid = git_reference_target(ref);
                if (oid) {
                    char hash[8];
                    git_oid_tostr(hash, sizeof(hash), oid);
                    string hashStr(hash, 7);
                    string name = git_reference_shorthand(ref);
                    refs[hashStr].push_back(name);
                }
                git_reference_free(ref);
            }
        }
        git_strarray_free(&ref_list);
    }

    // Обход коммитов
    git_oid oid;
    int count = 0;
    while (!git_revwalk_next(&oid, walk)) {
        git_commit* commit;
        git_commit_lookup(&commit, repo, &oid);
        // фильтр по автору
        if (!authorFilter.empty()) {
            const git_signature* author = git_commit_author(commit);
            if (string(author->name).find(authorFilter) == string::npos) {
                git_commit_free(commit);
                continue;
            }
        }
        count++;
        bool last = (count == maxCount || maxCount == 0);
        printCommit(commit, "", last, refs, decorate, oneline);
        git_commit_free(commit);
        if (maxCount > 0 && count >= maxCount) break;
    }

    git_revwalk_free(walk);
    git_repository_free(repo);
    git_libgit2_shutdown();
    return 0;
}
