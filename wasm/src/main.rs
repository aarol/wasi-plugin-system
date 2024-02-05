use std::io::{self, BufRead, Write};

use plugin::request;
use prost::Message;
use syntect::{highlighting::ThemeSet, parsing::SyntaxSet};

pub mod plugin {
    include!(concat!(env!("OUT_DIR"), "/plugin.rs"));
}

fn main() {
    let stdin = io::stdin();
    let mut stdin = stdin.lock();
    let buffer = stdin.fill_buf().unwrap();

    let ps = SyntaxSet::load_defaults_newlines();
    let ts = ThemeSet::load_defaults();

    let request = plugin::Request::decode(buffer).expect("Failed to decode request");

    match request.req {
        Some(request::Req::SyntaxRequest(r)) => {
            let syntax = ps.find_syntax_by_extension(&r.language).unwrap();
            let html = syntect::html::highlighted_html_for_string(
                &r.code,
                &ps,
                syntax,
                &ts.themes["base16-ocean.dark"],
            )
            .unwrap();

            let res = plugin::SyntaxResponse { output: html };
            io::stdout().write(&res.encode_to_vec()).unwrap();
        }
        None => todo!(),
    }
}
