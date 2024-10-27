use std::fmt;
use std::collections::BTreeMap;
use std::rc::Rc;

type DbRef = usize;

struct AtomicGroup {
	class:  DbRef,
	groups: Vec<DbRef>,
	tag:    String
}

impl fmt::Debug for AtomicGroup {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "AG::{}", self.tag)
    }
}

pub fn atomic_groups() {
    println!("Hello from atomicGroups!");

    let clid = 1;
    let mut divs: Vec<Vec<DbRef>> = Vec::new();
    divs.push(vec![1,2]);
    divs.push(vec![3,4]);
    divs.push(vec![5,6]);
    divs.push(vec![7,8,9]);
    println!("{:?}", divs);
    let mut agi: Vec<Vec<DbRef>> = vec![Vec::new()];
    for d in divs.iter() {        
        let mut agix: Vec<Vec<DbRef>> = Vec::new();
        for ag in agi.iter() {
            //println!("ag: {:?}", ag);
            for g in d.iter() {
                let mut gx = ag.clone();
                gx.push(*g);
                agix.push(gx);
            }
        }
        agi = agix;
    }
    println!("agi: {:?}", agi);

    // Make AtomicGroups
    let mut aglist: Vec<Rc<AtomicGroup>> = Vec::new();
    for ag in agi.iter() {
        let mut glist: Vec<String> = Vec::new();
        for g in  ag.iter() {
            glist.push(format!("G{}", g));
        }
        let gstr = format!("{}#{}", clid, glist.join("~"));
        //println!("#ag {:?}: {:?}", ag, gstr);
        let ago: AtomicGroup = AtomicGroup{
            class:  1,
            groups: ag.clone(),
            tag:    gstr
        };
        aglist.push(Rc::new(ago));
    }
    println!("aglist: {:?}", aglist);

    let mut g2ags: BTreeMap<DbRef, Vec<Rc<AtomicGroup>>> = BTreeMap::new();
    let mut i = divs.len();
    let mut n = 1;
    while i > 0 {
        i -= 1;
        let mut a = 0;
        while a < aglist.len() {
            for g in divs[i].iter() {
                for _ in 0 .. n {
                    g2ags.entry(*g).or_default().push(aglist[a].clone());
                    a += 1;
                }
            }
        }
        n *= divs[i].len();
    }
    println!("\n****************************\n  g2ags: {:?}", g2ags);
    
}

/*

		//fmt.Printf("     ++> %+v\n", xg2ags)
		if len(divs) != 0 {
			fetinfo.atomicgroups[cl.Id] = aglist
			for g, agl := range g2ags {
				agls := []string{}
				for _, ag := range agl {
					agls = append(agls, ag.Tag)
				}
				//fmt.Printf("     ++ %s: %+v\n", fetinfo.ref2fet[g], agls)
				fetinfo.atomicgroups[g] = agl
			}
		} else {
			fetinfo.atomicgroups[cl.Id] = []AtomicGroup{}
		}
	}
	//fmt.Println("\n +++++++++++++++++++++++++++")
	//printAtomicGroups(fetinfo)
}

func printAtomicGroups(fetinfo *fetInfo) {
	for _, cl := range fetinfo.db.Classes {
		agls := []string{}
		for _, ag := range fetinfo.atomicgroups[cl.Id] {
			agls = append(agls, ag.Tag)
		}
		fmt.Printf("  ++ %s: %+v\n", fetinfo.ref2fet[cl.Id], agls)
		for _, div := range fetinfo.classdivisions[cl.Id] {
			for _, g := range div {
				agls := []string{}
				for _, ag := range fetinfo.atomicgroups[g] {
					agls = append(agls, ag.Tag)
				}
				fmt.Printf("    -- %s: %+v\n", fetinfo.ref2fet[g], agls)
			}
		}
	}
}
 */