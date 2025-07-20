export namespace dtos {
	
	export class Mod {
	    id: string;
	    profile_id: string;
	    name: string;
	    order: number;
	    is_separator: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Mod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.name = source["name"];
	        this.order = source["order"];
	        this.is_separator = source["is_separator"];
	    }
	}
	export class GroupedMod {
	    separator: string;
	    mods: Mod[];
	
	    static createFrom(source: any = {}) {
	        return new GroupedMod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.separator = source["separator"];
	        this.mods = this.convertValues(source["mods"], Mod);
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
	
	export class ModlistDTO {
	    id: string;
	    name: string;
	    author: string;
	    description: string;
	    image: string;
	    game_type: string;
	    version: string;
	    is_nsfw: boolean;
	    website: string;
	    readme: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ModlistDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.author = source["author"];
	        this.description = source["description"];
	        this.image = source["image"];
	        this.game_type = source["game_type"];
	        this.version = source["version"];
	        this.is_nsfw = source["is_nsfw"];
	        this.website = source["website"];
	        this.readme = source["readme"];
	        this.created_at = source["created_at"];
	    }
	}

}

export namespace models {
	
	export class Profile {
	    id: string;
	    modlist_id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.modlist_id = source["modlist_id"];
	        this.name = source["name"];
	    }
	}
	export class ProfileFile {
	    id: string;
	    profile_id: string;
	    name: string;
	    file_path: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.name = source["name"];
	        this.file_path = source["file_path"];
	    }
	}

}

