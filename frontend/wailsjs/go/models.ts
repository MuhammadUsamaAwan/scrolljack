export namespace dtos {
	
	export class ModDTO {
	    id: string;
	    profile_id: string;
	    name: string;
	    order: number;
	    mod_order: number;
	    is_active: boolean;
	    is_separator: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.name = source["name"];
	        this.order = source["order"];
	        this.mod_order = source["mod_order"];
	        this.is_active = source["is_active"];
	        this.is_separator = source["is_separator"];
	    }
	}
	export class GroupedModDTO {
	    separator: string;
	    mods: ModDTO[];
	
	    static createFrom(source: any = {}) {
	        return new GroupedModDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.separator = source["separator"];
	        this.mods = this.convertValues(source["mods"], ModDTO);
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
	export class ModArchiveDTO {
	    id: string;
	    hash: string;
	    type: string;
	    nexus_game_name?: string;
	    nexus_mod_id?: string;
	    nexus_file_id?: string;
	    direct_url?: string;
	    version?: string;
	    size?: number;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new ModArchiveDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.hash = source["hash"];
	        this.type = source["type"];
	        this.nexus_game_name = source["nexus_game_name"];
	        this.nexus_mod_id = source["nexus_mod_id"];
	        this.nexus_file_id = source["nexus_file_id"];
	        this.direct_url = source["direct_url"];
	        this.version = source["version"];
	        this.size = source["size"];
	        this.description = source["description"];
	    }
	}
	
	export class ModFileDTO {
	    id: string;
	    hash: string;
	    type: string;
	    path: string;
	    source_file_path?: string;
	    patch_file_path?: string;
	    bsa_files?: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ModFileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.hash = source["hash"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.source_file_path = source["source_file_path"];
	        this.patch_file_path = source["patch_file_path"];
	        this.bsa_files = source["bsa_files"];
	        this.size = source["size"];
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

