export interface AlbumTrack {
  title: string;
  artist?: string;
  lrc?: string;
}

export interface MusicAlbum {
  slug: string;
  title: string;
  artist: string;
  year: string;
  cover: string;
  tracks: AlbumTrack[];
}

export const musicAlbums: MusicAlbum[] = [
  {
    slug: 'tsuki-wo-aruiteiru',
    title: '月を歩いている',
    artist: 'n-buna',
    year: '2016',
    cover: '/media/music/covers/tsuki-wo-aruiteiru.webp',
    tracks: [
      { title: 'モノローグ' },
      { title: 'ルラ', artist: 'n-buna feat. 初音ミク' },
      { title: '三月と狼少年', artist: 'n-buna feat. miki' },
      { title: '歌う睡蓮', artist: 'n-buna feat. 初音ミク' },
      { title: '花降らし', artist: 'n-buna feat. 初音ミク' },
      { title: '落花' },
      { title: '泣いた振りをした', artist: 'n-buna feat. 初音ミク' },
      { title: '白ゆき', artist: 'n-buna feat. 初音ミク' },
      { title: 'ラプンツェル', artist: 'n-buna feat. 初音ミク' },
      { title: '落陽' },
      { title: '白ゆきの独白' },
      { title: 'セロ弾き群青', artist: 'n-buna feat. 初音ミク' },
      { title: 'それでもいいよ。', artist: 'n-buna feat. 初音ミク' },
      { title: 'かぐや', artist: 'n-buna feat. 初音ミク' },
      { title: 'エピローグ' },
      { title: 'カエルのはなし', artist: 'n-buna feat. 初音ミク' },
      { title: '白ゆき PianoRock Arrange' },
    ],
  },
  {
    slug: 'hana-to-mizuame-saishuu-densha',
    title: '花と水飴、最終電車',
    artist: 'n-buna',
    year: '2015',
    cover: '/media/music/covers/hana-to-mizuame-saishuu-densha.webp',
    tracks: [
      { title: 'もうじき夏が終わるから' },
      { title: '無人駅' },
      { title: '始発とカフカ' },
      { title: 'ウミユリ海底譚' },
      { title: '昼青' },
      { title: '拝啓、夏に溺れる' },
      { title: 'ヒグレギ' },
      { title: '透明エレジー' },
      { title: '夜祭前に' },
      { title: 'メリュー' },
      { title: '着火、カウントダウン' },
      { title: '敬具' },
      { title: 'ずっと空を見ていた' },
      { title: '夜明けと蛍' },
      { title: '花と水飴、最終電車' },
    ],
  },
  {
    slug: 'dakara-boku-wa-ongaku-wo-yameta',
    title: 'だから僕は音楽を辞めた',
    artist: 'ヨルシカ',
    year: '2019',
    cover: '/media/music/covers/dakara-boku-wa-ongaku-wo-yameta.webp',
    tracks: [
      { title: '410' },
      { title: '藍二乗' },
      { title: '八月、某、月明かり' },
      { title: '詩書きとコーヒー' },
      { title: '713' },
      { title: '踊ろうぜ' },
      { title: '六月は雨上がりの街を書く' },
      { title: '五月は花緑青の窓辺から' },
      { title: '夜紛い' },
      { title: '56' },
      { title: 'パレード' },
      { title: 'エルマ' },
      { title: '831' },
      { title: 'だから僕は音楽を辞めた' },
    ],
  },
  {
    slug: 'curtain-call',
    title: 'カーテンコールが止む前に',
    artist: 'n-buna',
    year: '2014',
    cover: '/media/music/covers/curtain-call.webp',
    tracks: [
      { title: '一人きりロックショー', artist: 'n-buna feat. 初音ミク・GUMI' },
      { title: 'スロイド', artist: 'n-buna feat. 初音ミク' },
      { title: '透明エレジー', artist: 'n-buna feat. GUMI' },
      { title: 'アイラ', artist: 'n-buna feat. miki' },
      { title: 'また雨が降ったら', artist: 'n-buna feat. 初音ミク' },
      { title: '七月、影法師、藍色、ロッカー', artist: 'n-buna feat. GUMI・miki' },
      { title: '夕立' },
      { title: '背景、夏に溺れる', artist: 'n-buna feat. 初音ミク' },
      { title: 'カーテンコールが止む前に', artist: 'n-buna feat. 初音ミク' },
      { title: 'ウミユリ海底譚', artist: 'n-buna feat. 初音ミク' },
      { title: 'ハイカラ色の', artist: 'n-buna feat. 初音ミク' },
      { title: '夜に染まるまで', artist: 'n-buna feat. 初音ミク' },
      { title: '劇場愛歌', artist: 'n-buna feat. miki' },
      { title: 'さよならワンダーノイズ', artist: 'n-buna feat. 初音ミク' },
      { title: 'ウミユリ海底譚 piano ver.' },
    ],
  },
  {
    slug: 'elma',
    title: 'エルマ',
    artist: 'ヨルシカ',
    year: '2019',
    cover: '/media/music/covers/elma.webp',
    tracks: [
      { title: '車窓' },
      { title: '憂一乗' },
      { title: '夕凪、某、花惑い' },
      { title: '雨とカプチーノ' },
      { title: '湖の街' },
      { title: '神様のダンス' },
      { title: '雨晴るる' },
      { title: '歩く' },
      { title: '心に穴が空いた' },
      { title: '森の教会' },
      { title: '声' },
      { title: 'エイミー' },
      { title: '海底、月明かり' },
      { title: 'ノーチラス' },
    ],
  },
  {
    slug: 'sousaku',
    title: '創作',
    artist: 'ヨルシカ',
    year: '2021',
    cover: '/media/music/covers/sousaku.webp',
    tracks: [
      { title: '強盗と花束' },
      { title: '春泥棒' },
      { title: '創作' },
      { title: '風を食む' },
      { title: '嘘月' },
    ],
  },
  {
    slug: 'tousaku',
    title: '盗作',
    artist: 'ヨルシカ',
    year: '2020',
    cover: '/media/music/covers/tousaku.webp',
    tracks: [
      { title: '音楽泥棒の自白' },
      { title: '昼鳶' },
      { title: '春ひさぎ' },
      { title: '爆弾魔 (Re-Recording)' },
      { title: '青年期、空き巣' },
      { title: 'レプリカント' },
      { title: '花人局' },
      { title: '朱夏期、音楽泥棒' },
      { title: '盗作' },
      { title: '思想犯' },
      { title: '逃亡' },
      { title: '幼年期、思い出の中' },
      { title: '夜行' },
      { title: '花に亡霊' },
    ],
  },
];

export const albumStartIndexes = Object.fromEntries(
  musicAlbums.map((album, albumIndex) => [
    album.slug,
    musicAlbums.slice(0, albumIndex).reduce((total, item) => total + item.tracks.length, 0),
  ]),
) as Record<string, number>;

export const playerTracks = musicAlbums.flatMap((album) =>
  album.tracks.map((track, index) => ({
    name: track.title,
    artist: track.artist ?? album.artist,
    url: `/media/music/${album.slug}/${String(index + 1).padStart(2, '0')}.mp3`,
    lrc: track.lrc ?? `/media/music/${album.slug}/${String(index + 1).padStart(2, '0')}.lrc`,
    cover: album.cover,
    album: album.title,
  })),
);
