use clap::Parser;
use std::env;

mod error;
mod git;

pub use self::error::{Error, Result};

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    /// mode you wish to run, i.e., scan (s) or report (r)
    #[arg(short, long, default_value_t = 's')]
    mode: char,
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.mode {
        's' | 'r' => (),
        _ => todo!(),
    }

    let wd = env::current_dir()?;
    let repo = git::find_repo_root(&wd)?;
    println!("{}", repo.display());
    Ok(())
}
