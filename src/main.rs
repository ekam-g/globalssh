use ptyprocess::PtyProcess;
use std::io::{self, BufRead, BufReader, Result, Write};
use std::process::Command;
use std::thread;

fn main() -> Result<()> {
    // spawn a cat process
    let mut process = PtyProcess::spawn(Command::new("zsh"))?;
    // create a communication stream
    let mut stream = process.get_raw_handle()?;
    let input_thread = thread::spawn(move || {
        let mut input = String::new();
        loop {
            if io::stdin().read_line(&mut input).is_ok() {
                break;
            }
        }
        input
    });
    let mut reader = BufReader::new(stream);
    let mut buf = String::new();
    loop {
        reader.read_line(&mut buf)?;
        println!("{buf}");
        if let Ok(e) = process.is_alive() {
            if !e {
                break;
            }
        }
    }
    Ok(())
}
