use ptyprocess::PtyProcess;
use std::io::{BufRead, BufReader, Result, Write};
use std::process::Command;
use std::thread::{self};
use std::time::Duration;

fn main() -> Result<()> {
    // spawn a cat process
    let mut process = PtyProcess::spawn(Command::new("zsh"))?;

    // create a communication stream
    let mut stream = process.get_raw_handle()?;

    // send a message to process
    writeln!(stream, "neofetch")?;

    // read a line from the stream
    let mut reader = BufReader::new(stream);
    let mut buf = String::new();
    for _ in 0..100000 {
        reader.read_line(&mut buf)?;
        println!("{buf}");
    }


    // stop the process
    assert!(process.exit(true)?);

    Ok(())
}
