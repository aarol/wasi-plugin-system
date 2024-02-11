use std::{
    io::{self, BufRead, Read, Write},
    thread,
    time::{Duration, Instant},
};

use plugin::request;
use prost::Message;
use syntect::{highlighting::ThemeSet, parsing::SyntaxSet};

pub mod plugin {
    include!(concat!(env!("OUT_DIR"), "/wasi_plugin.rs"));
}

fn main() {
    let stdin = io::stdin();
    let mut stdin: io::StdinLock<'_> = stdin.lock();

    let mut buf = Vec::new();
    stdin.read_to_end(&mut buf);

    eprintln!("First input: {:?}", &buf);

    let mut buf = Vec::new();
    stdin.read_to_end(&mut buf);

    eprintln!("Second input: {:?}", &buf);

    return;
    // let request = plugin::Request::decode(read).expect("Failed to decode request");

    // match request.req {
    //     Some(request::Req::SyntaxRequest(r)) => {
    //         let now = Instant::now();
    //         let ps: SyntaxSet = SyntaxSet::load_defaults_newlines();
    //         let ts = ThemeSet::load_defaults();
    //         let syntax = ps.find_syntax_by_token(&r.language).unwrap();
    //         let html = syntect::html::highlighted_html_for_string(
    //             &r.code,
    //             &ps,
    //             syntax,
    //             &ts.themes["base16-ocean.dark"],
    //         )
    //         .unwrap();

    //         let res = plugin::SyntaxResponse { output: html };
    //         io::stdout().write(&res.encode_to_vec()).unwrap();

    //         write!(io::stderr(), "Elapsed: {}", now.elapsed().as_micros());
    //     }
    //     _ => {}
    // }
}
