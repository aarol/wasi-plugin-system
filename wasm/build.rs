use std::io::Result;

fn main() -> Result<()> {
    prost_build::compile_protos(&["../proto/plugin.proto"], &["../proto"])?;
    Ok(())
}
