use std::{thread, time::Duration, io::{self, BufReader}};


pub mod pty;



fn main()  {
    let tty = pty::PtyTerm::new("zsh").unwrap();
    thread::spawn(move ||  {
        let term = console::Term::stdout();
        loop {
            let value = console::Term::read_char(&term).unwrap();
            tty.input.send(value.to_string()).unwrap();
        }
    });
    loop {
        println!("{}",tty.output.recv().unwrap().unwrap());
    }
}

