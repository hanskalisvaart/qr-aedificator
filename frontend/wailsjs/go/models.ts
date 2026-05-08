export namespace main {
	
	export class BatchResult {
	    index: number;
	    image: string;
	    content: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.image = source["image"];
	        this.content = source["content"];
	        this.error = source["error"];
	    }
	}
	export class CSVParseResult {
	    headers: string[];
	    rows: any[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CSVParseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.headers = source["headers"];
	        this.rows = source["rows"];
	        this.error = source["error"];
	    }
	}
	export class HistoryEntry {
	    id: number;
	    type: string;
	    fields: string;
	    content: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.fields = source["fields"];
	        this.content = source["content"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class QRRequest {
	    type: string;
	    fields: Record<string, string>;
	    ecl: string;
	    size: number;
	    fgColor: string;
	    bgColor: string;
	
	    static createFrom(source: any = {}) {
	        return new QRRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.fields = source["fields"];
	        this.ecl = source["ecl"];
	        this.size = source["size"];
	        this.fgColor = source["fgColor"];
	        this.bgColor = source["bgColor"];
	    }
	}
	export class QRResponse {
	    image: string;
	    content: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new QRResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	        this.content = source["content"];
	        this.error = source["error"];
	    }
	}

}

