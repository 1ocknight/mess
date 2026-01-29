import defaultAvatar from '../../public/vite.svg';

export default function ProfileCard({ profile, onClick }) {
  if (!profile) return null;

  return (
    <div
      onClick={onClick}
      style={{
        display: 'flex',
        alignItems: 'center',
        cursor: 'pointer',
        padding: '12px 16px',
        borderRadius: 20,
        width: 220,
        background: 'linear-gradient(145deg, #7e57c2, #9c27b0)',
        color: '#fff',
        boxShadow: '0 6px 12px rgba(0,0,0,0.2)',
        transition: 'transform 0.2s, box-shadow 0.2s',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-4px)';
        e.currentTarget.style.boxShadow = '0 10px 20px rgba(0,0,0,0.3)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = '0 6px 12px rgba(0,0,0,0.2)';
      }}
    >
      <img
        src={profile.avatar_url || defaultAvatar}
        alt="Avatar"
        style={{
          width: 50,
          height: 50,
          borderRadius: '50%',
          marginRight: 12,
          border: '2px solid #d1b3ff',
          transition: 'transform 0.2s',
        }}
      />
      <span style={{ fontWeight: 'bold', fontSize: 16 }}>{profile.alias}</span>
    </div>
  );
}
