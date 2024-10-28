pub mod atomicgroups;
pub mod readw365;
use std::env;
use std::process::ExitCode;


fn main() -> ExitCode {
    println!("Hello from testag!");
    let path = env::current_dir().expect("Fail");
    println!("The current directory is {}", path.display());
    
    let inpath = "../internal/_testdata/test1r.json".to_string();
    let w365data = match readw365::read_w365(&inpath) {
        Ok(d) => d,
        Err(err) => {
            eprintln!("Error: {}:\n  \"{}\"", err, &inpath);
            return ExitCode::from(1);
        },
    };
    //println!("{:#?}", w365data);
    println!("{:#?}", w365data.Classes);

    //atomicgroups::atomic_groups();

    ExitCode::from(0)
}
