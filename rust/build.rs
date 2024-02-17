use std::io::Result;

fn main() -> Result<()> {
    prost_build::compile_protos(&["../plugin.proto"], &["../"])?;
    println!("cargo:rerun-if-changed=../plugin.proto");
    Ok(())
}
 