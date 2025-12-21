export namespace handler {
	
	export class ShortcutConfig {
	    key: string;
	    targetDir: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new ShortcutConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.targetDir = source["targetDir"];
	        this.label = source["label"];
	    }
	}

}

