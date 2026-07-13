use std::process::Command;

#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;

pub fn create_command<S: AsRef<std::ffi::OsStr>>(cmd: S) -> Command {
    let mut command = Command::new(cmd);
    
    #[cfg(target_os = "windows")]
    {
        // CREATE_NO_WINDOW = 0x08000000
        command.creation_flags(0x08000000);
    }
    
    command
}
