use std::sync::mpsc::{Receiver, Sender};

use anyhow::Error;
use portable_pty::{CommandBuilder, NativePtySystem, PtySize, PtySystem};

use std::{
    io::{BufRead, BufReader},
    sync::mpsc::channel,
};
#[derive(Debug)]
pub struct PtyTerm {
    input: Sender<String>,
    output: Receiver<Result<String,std::io::Error>>,
}

impl PtyTerm {
    pub fn new() -> Result<Self, Error> {
        let pty_system = NativePtySystem::default();
        //todo: change size system soon!
        let pair = pty_system
            .openpty(PtySize {
                rows: 24,
                cols: 80,
                pixel_width: 0,
                pixel_height: 0,
            })?;

        // Spawn a shell into the pty
        let cmd = CommandBuilder::new("zsh");
        let mut child = pair.slave.spawn_command(cmd)?;
        // Release any handles owned by the slave: we don't need it now
        // that we've spawned the child.
        drop(pair.slave);
        // Read the output in another thread.
        // This is important because it is easy to encounter a situation
        // where read/write buffers fill and block either your process
        // or the spawned process.
        let (tx, output) = channel();
        let reader = pair.master.try_clone_reader().unwrap();
        std::thread::spawn(move || {
            // Consume the output from the child
            let reader = BufReader::new(reader);
            for line in reader.lines() {
                tx.send(line);
            }
        });
        return Ok(PtyTerm{
            input: todo!(),
            output,
        });
    }
}
