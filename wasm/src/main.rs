use std::{
    io::{self, BufRead, Write},
    thread,
    time::Duration,
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
    let stdout = io::stdout();
    let mut stdout = stdout.lock();

    let info = plugin::PluginInfo {
        name: "".into(),
        events: vec!["syntax-highlight".to_string()],
    };

    stdout
        .write_all(&info.encode_length_delimited_to_vec())
        .expect("Write stdout");
    // let read = stdin.fill_buf().expect("Read stdin");

    // plugin::Request::decode_length_delimited(read).expect("Decode stdin");
    return;
    let ps = SyntaxSet::load_defaults_newlines();
    let ts = ThemeSet::load_defaults();

    loop {
        let mut buf = Vec::new();

        let request = plugin::Request::decode(buf.as_slice()).expect("Failed to decode request");

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
            _ => {}
        }
        io::stdout().write(b"\n").unwrap();
    }
}
