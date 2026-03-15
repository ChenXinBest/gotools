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
	export class ExportSettings {
	    export_tool: string;
	    export_path: string;
	    last_connection_id: string;
	    last_databases: string[];
	    last_database: string;
	    last_tables: string[];
	    threads: number;
	    skip_definer: boolean;
	    skip_binlog: boolean;
	    compression: string;
	    chunk_size: string;
	    export_scope: string;
	    include_schemas: string;
	    exclude_schemas: string;
	    include_tables: string;
	    exclude_tables: string;
	    overwrite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.export_tool = source["export_tool"];
	        this.export_path = source["export_path"];
	        this.last_connection_id = source["last_connection_id"];
	        this.last_databases = source["last_databases"];
	        this.last_database = source["last_database"];
	        this.last_tables = source["last_tables"];
	        this.threads = source["threads"];
	        this.skip_definer = source["skip_definer"];
	        this.skip_binlog = source["skip_binlog"];
	        this.compression = source["compression"];
	        this.chunk_size = source["chunk_size"];
	        this.export_scope = source["export_scope"];
	        this.include_schemas = source["include_schemas"];
	        this.exclude_schemas = source["exclude_schemas"];
	        this.include_tables = source["include_tables"];
	        this.exclude_tables = source["exclude_tables"];
	        this.overwrite = source["overwrite"];
	    }
	}

}

export namespace main {
	
	export class CheckImportConflictsRequest {
	    connection_id: string;
	    input_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckImportConflictsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.input_dir = source["input_dir"];
	    }
	}
	export class DropConflictingTablesRequest {
	    connection_id: string;
	    conflicts: tools.ImportConflict[];
	
	    static createFrom(source: any = {}) {
	        return new DropConflictingTablesRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.conflicts = this.convertValues(source["conflicts"], tools.ImportConflict);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExportRequest {
	    connection_id: string;
	    databases: string[];
	    database: string;
	    tables: string[];
	    output_dir: string;
	    threads: number;
	    compression: string;
	    chunk_size: string;
	    skip_definer: boolean;
	    skip_binlog: boolean;
	    include_schemas: string[];
	    exclude_schemas: string[];
	    include_tables: string[];
	    exclude_tables: string[];
	    overwrite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.databases = source["databases"];
	        this.database = source["database"];
	        this.tables = source["tables"];
	        this.output_dir = source["output_dir"];
	        this.threads = source["threads"];
	        this.compression = source["compression"];
	        this.chunk_size = source["chunk_size"];
	        this.skip_definer = source["skip_definer"];
	        this.skip_binlog = source["skip_binlog"];
	        this.include_schemas = source["include_schemas"];
	        this.exclude_schemas = source["exclude_schemas"];
	        this.include_tables = source["include_tables"];
	        this.exclude_tables = source["exclude_tables"];
	        this.overwrite = source["overwrite"];
	    }
	}
	export class ExportResponse {
	    success: boolean;
	    message: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.path = source["path"];
	    }
	}
	export class ImportRequest {
	    connection_id: string;
	    database: string;
	    input_dir: string;
	    threads: number;
	    schema: string;
	    include_schemas: string[];
	    exclude_schemas: string[];
	    include_tables: string[];
	    exclude_tables: string[];
	    reset_progress: boolean;
	    wait_timeout: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.database = source["database"];
	        this.input_dir = source["input_dir"];
	        this.threads = source["threads"];
	        this.schema = source["schema"];
	        this.include_schemas = source["include_schemas"];
	        this.exclude_schemas = source["exclude_schemas"];
	        this.include_tables = source["include_tables"];
	        this.exclude_tables = source["exclude_tables"];
	        this.reset_progress = source["reset_progress"];
	        this.wait_timeout = source["wait_timeout"];
	    }
	}
	export class ImportResponse {
	    success: boolean;
	    message: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.path = source["path"];
	    }
	}
	export class MySQLDumpExportRequest {
	    connection_id: string;
	    databases: string[];
	    database: string;
	    tables: string[];
	    output_dir: string;
	    compression: string;
	    single_transaction: boolean;
	    routines: boolean;
	    events: boolean;
	    overwrite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MySQLDumpExportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.databases = source["databases"];
	        this.database = source["database"];
	        this.tables = source["tables"];
	        this.output_dir = source["output_dir"];
	        this.compression = source["compression"];
	        this.single_transaction = source["single_transaction"];
	        this.routines = source["routines"];
	        this.events = source["events"];
	        this.overwrite = source["overwrite"];
	    }
	}
	export class MySQLDumpImportRequest {
	    connection_id: string;
	    input_file: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new MySQLDumpImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_id = source["connection_id"];
	        this.input_file = source["input_file"];
	        this.database = source["database"];
	    }
	}

}

export namespace tools {
	
	export class ImportConflict {
	    schema: string;
	    tables: string[];
	    views: string[];
	    events: string[];
	    functions: string[];
	    procedures: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImportConflict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.tables = source["tables"];
	        this.views = source["views"];
	        this.events = source["events"];
	        this.functions = source["functions"];
	        this.procedures = source["procedures"];
	    }
	}
	export class ImportConflictCheckResult {
	    has_conflicts: boolean;
	    conflicts: ImportConflict[];
	
	    static createFrom(source: any = {}) {
	        return new ImportConflictCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_conflicts = source["has_conflicts"];
	        this.conflicts = this.convertValues(source["conflicts"], ImportConflict);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
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

export namespace version {
	
	export class Info {
	    version: string;
	    build_time: string;
	    git_commit: string;
	    git_branch: string;
	    go_version: string;
	    platform: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.build_time = source["build_time"];
	        this.git_commit = source["git_commit"];
	        this.git_branch = source["git_branch"];
	        this.go_version = source["go_version"];
	        this.platform = source["platform"];
	    }
	}

}

