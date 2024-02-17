use std::{
    error::Error,
    io::{self, Read, Write},
};
use syntect::{highlighting::ThemeSet, parsing::SyntaxSet};

use plugin::request::Req;
use prost::Message;

pub mod plugin {
    include!(concat!(env!("OUT_DIR"), "/wasi_plugin.rs"));
}

fn main() {
    if let Err(e) = run() {
        eprintln!("Encountered error: {e:?}")
    }
}

fn run() -> Result<(), Box<dyn Error>> {
    let stdin = io::stdin();
    let mut stdin: io::StdinLock<'_> = stdin.lock();

    let mut buf = Vec::new();
    stdin.read_to_end(&mut buf)?;

    let request = plugin::Request::decode(buf.as_slice())?;

    match request.req {
        Some(Req::SyntaxRequest(r)) => {
            let ps: SyntaxSet = SyntaxSet::load_defaults_newlines();
            let ts = ThemeSet::load_defaults();
            let syntax = ps
                .find_syntax_by_token(&r.language)
                .ok_or("Requested language syntax not found")?;
            let html = syntect::html::highlighted_html_for_string(
                &r.code,
                &ps,
                syntax,
                &ts.themes["base16-ocean.dark"],
            )?;

            let res = plugin::SyntaxResponse { output: html };
            io::stdout().write_all(&res.encode_to_vec())?;
        }
        Some(Req::VersionRequest(_)) => {
            let res = plugin::VersionResponse {
                version: "1.0.0".to_owned(),
            };
            io::stdout().write_all(&res.encode_to_vec())?;
        }
        _ => {}
    }
    Ok(())
}
