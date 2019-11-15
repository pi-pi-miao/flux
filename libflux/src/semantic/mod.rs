mod analyze;
pub use analyze::{analyze_with, Analyzer};

mod env;
pub mod fresh;
mod infer;

// TODO(jsternberg): Once more work is done on the infer methods,
// this should be removed and the warnings fixed.
#[allow(warnings)]
pub mod nodes;

mod sub;
mod types;
pub mod walk;

#[cfg(test)]
mod parser;

#[cfg(test)]
mod tests;

use crate::ast;
use crate::parser::parse_string;
use crate::semantic::analyze::Result;

pub fn analyze_source(source: &str, f: &mut fresh::Fresher) -> Result<nodes::Package> {
    let file = parse_string("", source);
    let errs = ast::check::check(ast::walk::Node::File(&file));
    if errs.len() > 0 {
        return Err(format!("got errors on parsing: {:?}", errs));
    }
    let ast_pkg = ast::Package {
        base: file.base.clone(),
        path: "".to_string(),
        package: "main".to_string(),
        files: vec![file],
    };
    analyze_with(ast_pkg, f)
}
