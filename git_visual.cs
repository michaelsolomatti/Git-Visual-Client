// git_visual.cs
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using LibGit2Sharp;

class GitVisual
{
    private string repoPath;
    private bool graph;
    private bool all;
    private bool decorate;
    private bool oneline;
    private string author;
    private string since;
    private string until;
    private int maxCount;
    private bool color;
    private Repository repo;

    public GitVisual(string path, bool graph, bool all, bool decorate, bool oneline,
                     string author, string since, string until, int maxCount, bool color)
    {
        this.repoPath = path;
        this.graph = graph;
        this.all = all;
        this.decorate = decorate;
        this.oneline = oneline;
        this.author = author;
        this.since = since;
        this.until = until;
        this.maxCount = maxCount;
        this.color = color;
    }

    public void OpenRepo()
    {
        repo = new Repository(repoPath);
    }

    private string Colorize(string text, string colorCode)
    {
        if (!color) return text;
        var codes = new Dictionary<string, string>
        {
            {"red", "\x1b[91m"},
            {"green", "\x1b[92m"},
            {"yellow", "\x1b[93m"},
            {"blue", "\x1b[94m"},
            {"cyan", "\x1b[96m"},
            {"reset", "\x1b[0m"}
        };
        return codes[colorCode] + text + codes["reset"];
    }

    private IEnumerable<Commit> GetCommits()
    {
        var filter = new CommitFilter
        {
            SortBy = CommitSortStrategies.Time
        };
        if (all)
            filter.IncludeReachableFrom = repo.Refs;
        else
            filter.IncludeReachableFrom = repo.Head;
        if (!string.IsNullOrEmpty(author))
            filter.Author = author;
        // since/until можно добавить через фильтрацию после
        var commits = repo.Commits.QueryBy(filter);
        if (maxCount > 0)
            commits = commits.Take(maxCount);
        return commits;
    }

    private Dictionary<string, List<string>> GetRefs()
    {
        var dict = new Dictionary<string, List<string>>();
        foreach (var branch in repo.Branches)
        {
            if (branch.Tip != null)
            {
                string hash = branch.Tip.Sha.Substring(0, 7);
                if (!dict.ContainsKey(hash)) dict[hash] = new List<string>();
                dict[hash].Add(branch.FriendlyName);
            }
        }
        foreach (var tag in repo.Tags)
        {
            var target = tag.Target as Commit;
            if (target != null)
            {
                string hash = target.Sha.Substring(0, 7);
                if (!dict.ContainsKey(hash)) dict[hash] = new List<string>();
                dict[hash].Add(tag.FriendlyName);
            }
        }
        // HEAD
        if (repo.Head != null && repo.Head.Tip != null)
        {
            string hash = repo.Head.Tip.Sha.Substring(0, 7);
            if (!dict.ContainsKey(hash)) dict[hash] = new List<string>();
            dict[hash].Add("HEAD");
        }
        return dict;
    }

    public void ShowGraph()
    {
        var commits = GetCommits();
        var refs = GetRefs();
        int count = commits.Count();
        int i = 0;
        foreach (var c in commits)
        {
            bool last = i == count - 1;
            string connector = last ? "└── " : "├── ";
            string hashShort = c.Sha.Substring(0, 7);
            string refStr = "";
            if (decorate && refs.ContainsKey(hashShort))
                refStr = " (" + string.Join(", ", refs[hashShort]) + ")";
            if (oneline)
            {
                Console.WriteLine($"{connector}{Colorize(hashShort, "yellow")}{Colorize(refStr, "blue")} {Colorize(c.MessageShort, "reset")}");
            }
            else
            {
                Console.WriteLine($"{connector}{Colorize(hashShort, "yellow")}{Colorize(refStr, "blue")}");
                Console.WriteLine($"    {Colorize("Author:", "cyan")} {c.Author.Name} <{c.Author.Email}>");
                Console.WriteLine($"    {Colorize("Date:", "green")} {c.Author.When}");
                Console.WriteLine($"    {Colorize(c.MessageShort, "reset")}\n");
            }
            i++;
        }
    }

    static void Main(string[] args)
    {
        var path = ".";
        var graph = true;
        var all = false;
        var decorate = true;
        var oneline = false;
        string author = null;
        string since = null;
        string until = null;
        int maxCount = 0;
        bool color = true;

        for (int i = 0; i < args.Length; i++)
        {
            switch (args[i])
            {
                case "--path": path = args[++i]; break;
                case "--graph": graph = bool.Parse(args[++i]); break;
                case "--all": all = true; break;
                case "--decorate": decorate = bool.Parse(args[++i]); break;
                case "--oneline": oneline = true; break;
                case "--author": author = args[++i]; break;
                case "--since": since = args[++i]; break;
                case "--until": until = args[++i]; break;
                case "-n": maxCount = int.Parse(args[++i]); break;
                case "--color": color = bool.Parse(args[++i]); break;
                case "--help": Console.WriteLine("Usage: ..."); return;
            }
        }

        var visual = new GitVisual(path, graph, all, decorate, oneline, author, since, until, maxCount, color);
        try
        {
            visual.OpenRepo();
            visual.ShowGraph();
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"Error: {ex.Message}");
            Environment.Exit(1);
        }
    }
}
