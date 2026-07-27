// git_visual.rs
use std::env;
use std::path::Path;
use std::process;
use git2::{Repository, Commit, ObjectType, BranchType, Reference, Oid};
use structopt::StructOpt;
use chrono::{DateTime, Local};

#[derive(Debug, StructOpt)]
#[structopt(name = "git_visual")]
struct Opt {
    #[structopt(long, default_value = ".")]
    path: String,

    #[structopt(long, default_value = "true")]
    graph: bool,

    #[structopt(long)]
    all: bool,

    #[structopt(long, default_value = "true")]
    decorate: bool,

    #[structopt(long)]
    oneline: bool,

    #[structopt(long)]
    author: Option<String>,

    #[structopt(long)]
    since: Option<String>,

    #[structopt(long)]
    until: Option<String>,

    #[structopt(short = "n", long)]
    max_count: Option<usize>,

    #[structopt(long, default_value = "true")]
    color: bool,
}

fn colorize(text: &str, color: &str, enabled: bool) -> String {
    if !enabled {
        return text.to_string();
    }
    let code = match color {
        "red" => "\x1b[91m",
        "green" => "\x1b[92m",
        "yellow" => "\x1b[93m",
        "blue" => "\x1b[94m",
        "cyan" => "\x1b[96m",
        _ => "\x1b[0m",
    };
    format!("{}{}\x1b[0m", code, text)
}

fn main() {
    let opt = Opt::from_args();
    let repo = match Repository::open(&opt.path) {
        Ok(r) => r,
        Err(e) => {
            eprintln!("Error opening repo: {}", e);
            process::exit(1);
        }
    };

    let mut revwalk = repo.revwalk().unwrap();
    if opt.all {
        revwalk.push_glob("refs/*").unwrap();
    } else {
        let head = repo.head().unwrap();
        revwalk.push(head.target().unwrap()).unwrap();
    }
    revwalk.set_sorting(git2::Sort::TIME).unwrap();

    let commits: Vec<Commit> = revwalk
        .take(opt.max_count.unwrap_or(usize::MAX))
        .map(|oid| repo.find_commit(oid.unwrap()).unwrap())
        .collect();

    // Собираем референсы для декора
    let mut refs_map: std::collections::HashMap<String, Vec<String>> = std::collections::HashMap::new();
    if opt.decorate {
        // Ветки
        let branches = repo.branches(Some(BranchType::Local)).unwrap();
        for branch in branches {
            let (branch, _) = branch.unwrap();
            let name = branch.name().unwrap().unwrap().to_string();
            let oid = branch.get().target().unwrap();
            let short = oid.to_string()[0..7].to_string();
            refs_map.entry(short).or_insert(Vec::new()).push(name);
        }
        // Теги
        let tags = repo.tag_names(None).unwrap();
        for tag_name in tags.iter().flatten() {
            if let Ok(oid) = repo.refname_to_id(&format!("refs/tags/{}", tag_name)) {
                let short = oid.to_string()[0..7].to_string();
                refs_map.entry(short).or_insert(Vec::new()).push(tag_name.to_string());
            }
        }
        // HEAD
        if let Ok(head) = repo.head() {
            if let Some(oid) = head.target() {
                let short = oid.to_string()[0..7].to_string();
                let name = if head.is_branch() {
                    format!("HEAD -> {}", head.shorthand().unwrap_or(""))
                } else {
                    "HEAD".to_string()
                };
                refs_map.entry(short).or_insert(Vec::new()).push(name);
            }
        }
    }

    let color = opt.color;

    for (i, commit) in commits.iter().enumerate() {
        let last = i == commits.len() - 1;
        let connector = if last { "└── " } else { "├── " };
        let hash = commit.id().to_string();
        let short_hash = &hash[0..7];
        let author = commit.author();
        let author_name = author.name().unwrap_or("unknown");
        let email = author.email().unwrap_or("");
        let date = DateTime::<Local>::from(commit.time().seconds()).format("%Y-%m-%d %H:%M").to_string();
        let msg = commit.message().unwrap_or("").lines().next().unwrap_or("");

        let ref_str = if opt.decorate {
            if let Some(refs) = refs_map.get(short_hash) {
                format!(" ({})", refs.join(", "))
            } else {
                String::new()
            }
        } else {
            String::new()
        };

        if opt.oneline {
            println!("{}{}{}{}{}",
                connector,
                colorize(short_hash, "yellow", color),
                colorize(&ref_str, "blue", color),
                " ",
                colorize(msg, "reset", color));
        } else {
            println!("{}{}{}{}",
                connector,
                colorize(short_hash, "yellow", color),
                colorize(&ref_str, "blue", color));
            println!("{}    {} {} <{}>",
                "",
                colorize("Author:", "cyan", color),
                author_name,
                email);
            println!("{}    {} {}",
                "",
                colorize("Date:", "green", color),
                date);
            println!("{}    {}\n",
                "",
                colorize(msg, "reset", color));
        }
    }
}
