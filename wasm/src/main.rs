use std::io::{self, BufRead, Cursor, Write};

use prost::Message;

pub mod plugin {
    include!(concat!(env!("OUT_DIR"), "/plugin.rs"));
}

fn main() {
    let stdin = io::stdin();
    let mut stdin = stdin.lock();
    let buffer = stdin.fill_buf().unwrap();
    match plugin::Request::decode(buffer) {
        Ok(req) => {
            let res = plugin::Response { output: req.input };
            io::stdout().write(&res.encode_to_vec()).unwrap();
        }
        Err(err) => panic!("{}", err),
    }
}
