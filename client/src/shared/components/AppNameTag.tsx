interface tag {
    tagline?: string
}

const AppNameTag = ({ tagline }: tag ) => (
  <>
    <div className='app-name'> MatchMe </div>
    <div className='summary'>
      <p className="catch-phrase">
        { tagline ?? 'Where real connections begin..'}
      </p>
    </div>
  </>
)

export default AppNameTag