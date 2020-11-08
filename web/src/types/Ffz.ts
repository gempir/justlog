export interface FfzChannelEmotesResponse {
    // room: Room;
    sets: Sets;
}

export interface FfzGlobalEmotesResponse {
    sets: Sets;
}

// export interface Room {
//     _id:             number;
//     css:             null;
//     display_name:    string;
//     id:              string;
//     is_group:        boolean;
//     mod_urls:        null;
//     moderator_badge: null;
//     set:             number;
//     twitch_id:       number;
//     user_badges:     UserBadges;
// }

// export interface UserBadges {
// }

export interface Sets {
    [key: string]: EmoteSet;
}

export interface EmoteSet {
    // _type:       number;
    // css:         null;
    // description: null;
    emoticons:   Emoticon[];
    // icon:        null;
    // id:          number;
    // title:       string;
}

export interface Emoticon {
    // css:      null;
    // height:   number;
    // hidden:   boolean;
    id:       number;
    // margins:  null;
    // modifier: boolean;
    name:     string;
    // offset:   null;
    // owner:    Owner;
    // public:   boolean;
    urls:     { [key: string]: string };
    // width:    number;
}

export interface Owner {
    _id:          number;
    display_name: string;
    name:         string;
}