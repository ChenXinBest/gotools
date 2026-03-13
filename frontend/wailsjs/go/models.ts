export namespace config {
	
	export class DatabaseConnection {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    user: string;
	    password: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.database = source["database"];
	    }
	}

}

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

