export namespace cpolar {
	
	export class Tunnel {
	    id: string;
	    name: string;
	    public_url: string;
	    proto: string;
	    addr: string;
	    create_datetime: string;
	
	    static createFrom(source: any = {}) {
	        return new Tunnel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.public_url = source["public_url"];
	        this.proto = source["proto"];
	        this.addr = source["addr"];
	        this.create_datetime = source["create_datetime"];
	    }
	}
	export class User {
	    Tunnels: Tunnel[];
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Tunnels = this.convertValues(source["Tunnels"], Tunnel);
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

}

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

