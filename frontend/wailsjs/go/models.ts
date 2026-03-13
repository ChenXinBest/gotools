export namespace tools {
	
	export class ProcessInfo {
	    PID: number;
	    Name: string;
	    Cmdline: string;
	    CPUPercent: number;
	    MemoryMB: number;
	    ListenAddr: string;
	    ListenPort: number;
	    Status: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PID = source["PID"];
	        this.Name = source["Name"];
	        this.Cmdline = source["Cmdline"];
	        this.CPUPercent = source["CPUPercent"];
	        this.MemoryMB = source["MemoryMB"];
	        this.ListenAddr = source["ListenAddr"];
	        this.ListenPort = source["ListenPort"];
	        this.Status = source["Status"];
	    }
	}

}

