use std::sync::mpsc::{Receiver, Sender};

use anyhow::Error;
use portable_pty::{CommandBuilder, NativePtySystem, PtySize, PtySystem};

use std::{
    io::{BufRead, BufReader},
    sync::mpsc::channel,
};
#[derive(Debug)]
pub struct PtyTerm {
    pub input: Sender<String>,
    pub output: Receiver<Result<String,std::io::Error>>,
}

impl PtyTerm {
    pub fn new(term : &str) -> Result<Self, Error> {
        let pty_system = NativePtySystem::default();
        //todo: change size system soon!
        let pair = pty_system
            .openpty(PtySize {
                rows: 24,
                cols: 80,
                pixel_width: 0,
                pixel_height: 0,
            })?;

        let cmd = CommandBuilder::new(term);
        pair.slave.spawn_command(cmd)?;
        // Release any handles owned by the slave: we don't need it now
        // that we've spawned the child.
        drop(pair.slave);
        let (tx, output) = channel();
        let reader = pair.master.try_clone_reader().unwrap();
        std::thread::spawn(move || {
            // Consume the output from the child
            let reader = BufReader::new(reader);
            for line in reader.lines() {
                tx.send(line).expect("Unexpected Thread Closure, Crital Failure");
            }
        });
        let mut writer = pair.master.take_writer().unwrap();
        let (input, rx) = channel();
        if cfg!(target_os = "macos") {
            // macOS quirk: the child and reader must be started and
            // allowed a brief grace period to run before we allow
            // the writer to drop. Otherwise, the data we send to
            // the kernel to trigger EOF is interleaved with the
            // data read by the reader! WTF!?
            std::thread::sleep(std::time::Duration::from_millis(20));
        }
        std::thread::spawn(move || loop {
            let data = rx.recv().expect("Critial Thread Error");
            if let  Err(e) = writeln!(writer, "{data}") {
                log::error!("Failed to write due to: {e}");
            }
        });
        Ok(PtyTerm{
            input,
            output,
        })
    }
}
