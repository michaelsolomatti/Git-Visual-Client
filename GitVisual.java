// GitVisual.java
import org.eclipse.jgit.api.Git;
import org.eclipse.jgit.api.LogCommand;
import org.eclipse.jgit.lib.Repository;
import org.eclipse.jgit.lib.Ref;
import org.eclipse.jgit.revwalk.RevCommit;
import org.eclipse.jgit.revwalk.RevWalk;
import org.eclipse.jgit.storage.file.FileRepositoryBuilder;
import java.io.File;
import java.text.SimpleDateFormat;
import java.util.*;

public class GitVisual {
    private static boolean color = true;
    private static boolean graph = true;
    private static boolean all = false;
    private static boolean decorate = true;
    private static boolean oneline = false;
    private static String author = null;
    private static int maxCount = 0;
    private static String repoPath = ".";

    public static void main(String[] args) throws Exception {
        parseArgs(args);
        Repository repo = new FileRepositoryBuilder()
                .setGitDir(new File(repoPath, ".git"))
                .build();
        try (Git git = new Git(repo)) {
            LogCommand log = git.log();
            if (all) log.all();
            // фильтры можно добавить через setSkip
            Iterable<RevCommit> commits = log.call();
            // собираем референсы
            Map<String, List<String>> refs = new HashMap<>();
            if (decorate) {
                for (Ref ref : repo.getAllRefs().values()) {
                    if (ref.getObjectId() != null) {
                        String hash = ref.getObjectId().getName().substring(0, 7);
                        refs.computeIfAbsent(hash, k -> new ArrayList<>()).add(ref.getName());
                    }
                }
            }
            int count = 0;
            for (RevCommit c : commits) {
                if (maxCount > 0 && count >= maxCount) break;
                if (author != null && !c.getAuthorIdent().getName().contains(author)) continue;
                count++;
                boolean last = count == maxCount; // упрощённо
                String connector = (count == 1) ? "└── " : "├── ";
                // в реальности нужно определять последний, но для демонстрации
                String hash = c.getId().getName().substring(0, 7);
                String refStr = "";
                if (decorate && refs.containsKey(hash)) {
                    refStr = " (" + String.join(", ", refs.get(hash)) + ")";
                }
                if (oneline) {
                    System.out.println(connector + colorize(hash, "yellow") + colorize(refStr, "blue") + " " + colorize(c.getShortMessage(), "reset"));
                } else {
                    System.out.println(connector + colorize(hash, "yellow") + colorize(refStr, "blue"));
                    System.out.println("    " + colorize("Author:", "cyan") + " " + c.getAuthorIdent().getName() + " <" + c.getAuthorIdent().getEmailAddress() + ">");
                    SimpleDateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm");
                    System.out.println("    " + colorize("Date:", "green") + " " + df.format(c.getAuthorIdent().getWhen()));
                    System.out.println("    " + colorize(c.getShortMessage(), "reset") + "\n");
                }
            }
        }
    }

    private static String colorize(String text, String colorCode) {
        if (!color) return text;
        Map<String, String> codes = new HashMap<>();
        codes.put("red", "\033[91m");
        codes.put("green", "\033[92m");
        codes.put("yellow", "\033[93m");
        codes.put("blue", "\033[94m");
        codes.put("cyan", "\033[96m");
        codes.put("reset", "\033[0m");
        return codes.getOrDefault(colorCode, "") + text + codes.get("reset");
    }

    private static void parseArgs(String[] args) {
        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
                case "--path": repoPath = args[++i]; break;
                case "--graph": graph = Boolean.parseBoolean(args[++i]); break;
                case "--all": all = true; break;
                case "--decorate": decorate = Boolean.parseBoolean(args[++i]); break;
                case "--oneline": oneline = true; break;
                case "--author": author = args[++i]; break;
                case "-n": maxCount = Integer.parseInt(args[++i]); break;
                case "--color": color = Boolean.parseBoolean(args[++i]); break;
                case "--help": System.out.println("Usage: ..."); System.exit(0);
            }
        }
    }
}
