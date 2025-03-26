export namespace main {
	
	export class UserInfo {
	    uname: string;
	    uemail: string;
	    uphone: string;
	    upassword: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uname = source["uname"];
	        this.uemail = source["uemail"];
	        this.uphone = source["uphone"];
	        this.upassword = source["upassword"];
	    }
	}

}

