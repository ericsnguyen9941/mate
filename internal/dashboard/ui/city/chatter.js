// Small talk: lines going round Vietnamese social media that a mascot says
// now and then while it works, thinks, stands idle or sleeps (CON-128).
// They are for fun and say nothing about the task, so they come in their own
// yellow bubble, never the paper one real status lines use, and never in the
// tent's feed.

export const CHAT_STYLE = { background: '#FFF4C7', color: '#5A4A12', border: 'rgba(160,120,20,0.32)', dots: '#8A7432' };

export const LINES = {
  working: [
    'Đỉnh nóc, kịch trần, bay phấp phới!',
    'Gét gô!',
    'Code như chưa từng được code',
    'Hơi bị đỉnh rồi đó nha',
    'Ổn áp luôn',
    'Mãi đỉnh!',
    'Chấn động địa cầu',
    'Đâu phải dạng vừa đâu',
    'Khá khen cho thằng em',
    'Tới công chuyện luôn!',
    'Let him cook',
    'Nhìn thì dễ, làm mới thấy',
    'Xin vĩnh biệt cụ bug',
    'Chạy được là được',
    'It works on my machine',
    'Một chín một mười',
    'Còn cái nịt bug!',
    'Chill guy đang code',
    'Đỉnh của chóp',
    'Ét o ét, deadline!',
  ],
  thinking: [
    'Ủa alo?',
    'Hỏi chấm???',
    'Tôi là ai, đây là đâu?',
    'Căng đét',
    'Ủa gì kỳ vậy ta',
    'Ơ kìa...',
    'Khoan, để em suy nghĩ',
    'Cái này khó à nha',
    'Hông hiểu luôn á',
    'Sai sai ở đâu đó',
    'Nghĩ nát óc luôn',
    'Ối dồi ôi',
    'Bình tĩnh, tự tin, không cay',
    'Hỏi Mate phát cho chắc',
  ],
  idle: [
    'Chill chill thôi',
    'Cho xin cái task',
    'Ngồi chơi xơi nước',
    'Hết việc rồi hả?',
    'Ét o ét!',
    'Rảnh quá trời',
    'Tâm trạng: chờ review',
    'Mlem mlem',
    'Đời là bể khổ, task là bờ',
    'Ai rủ đi cà phê không?',
  ],
  completed: ['Xong rồi nha!', 'Ngon lành cành đào', 'Quá đã!', 'Đỉnh thật sự'],
  sleeping: ['5 phút nữa thôi...', 'Đừng gọi, đang ngủ đông', 'Zzz... mơ thấy merge'],
};

// How often, in seconds of the clock it runs on: the next line comes
// GAP[0] + random * GAP[1] after the last.
export const GAP = [12, 22];

/** A seeded random in [0, 1), so a tent rebuilt from the same history chats the same. */
export function seeded(text) {
  let a = 2166136261;
  for (const ch of String(text)) a = Math.imul(a ^ ch.codePointAt(0), 16777619);
  return () => {
    a = (a + 0x6d2b79f5) | 0;
    let t = Math.imul(a ^ (a >>> 15), 1 | a);
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

/** A line for a mascot in `state`, or null when the state has none. */
export function lineFor(state, rand = Math.random) {
  const list = LINES[state];
  return list ? list[Math.floor(rand() * list.length)] : null;
}
