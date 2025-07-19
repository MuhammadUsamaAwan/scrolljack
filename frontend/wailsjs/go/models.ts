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
	        this.created_at = source["created_at"];
	    }
	}

}

