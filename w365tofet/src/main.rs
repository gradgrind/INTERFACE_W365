pub mod atomicgroups;
pub mod readw365;
use std::env;


fn main() {
    println!("Hello from testag!");
    let path = env::current_dir().expect("Fail");
    println!("The current directory is {}", path.display());

    
    readw365::read_w365("../internal/_testdata/test1r.json".to_string());

    //atomicgroups::atomic_groups();
}