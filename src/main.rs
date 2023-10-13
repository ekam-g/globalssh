use std::{thread, time::Duration};


pub mod pty;



fn main()  {
    let tty = pty::PtyTerm::new("zsh").unwrap();
    thread::spawn(move || loop {
        tty.input.send("ls -a\n".to_owned()).unwrap();
        thread::sleep(Duration::from_secs(1));
    });
    loop {
        println!("{}",tty.output.recv().unwrap().unwrap());
    }
}
