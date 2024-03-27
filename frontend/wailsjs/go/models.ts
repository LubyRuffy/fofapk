export namespace main {
	
	export class Result {
	    data: any;
	    code: number;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.code = source["code"];
	        this.error = source["error"];
	    }
	}

}

