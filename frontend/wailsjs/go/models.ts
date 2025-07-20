export namespace dtos {
	
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

