use crate::{Error, Result};
use std::{cmp::Ordering, path::Path};

pub struct Repo {
    // Maximum repo name len is 100 characters for GHub, GLab & bitbucket
    pub repo_name: [u8; 100],
    // Maximum user name len seems to be 40 chars. Might need to fact check this
    pub repo_owner: [u8; 40],
    pub repo_root: Path,
}

pub fn find_repo_root(path: &Path) -> Result<&Path> {
    if path.cmp(Path::new("/")) == Ordering::Equal {
        return Err(Error::custom("failed to locate repo root"));
    }

    println!("Directory: {}", path.display());

    if path.join(".git").exists() {
        return Ok(path);
    }

    if let Some(parent) = path.parent() {
        find_repo_root(parent)
    } else {
        Err(Error::custom("failed to get parent directory"))
    }
}

#[cfg(test)]
mod tests {
    use std::env;

    use super::*;

    #[test]
    fn find_repo_root_ok() {
        assert!(find_repo_root(&env::current_dir().unwrap()).is_ok());
    }

    #[test]
    fn find_repo_root_fail() {
        assert!(find_repo_root(Path::new("/Home/User/project")).is_err());
    }
}
