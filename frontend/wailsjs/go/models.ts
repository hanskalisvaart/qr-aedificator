export namespace main {
	
	export class QRRequest {
	    type: string;
	    fields: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new QRRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.fields = source["fields"];
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

