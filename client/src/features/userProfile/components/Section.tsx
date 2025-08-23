type SectionProps = {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
  className?: string;
};

const Section = ({ title, subtitle, children, className }: SectionProps) => (
  <div className={`section ${className ?? ""}`}>
    <div className="section-header">
      <div>
        <div className="section-title">{title}</div>
        {subtitle && <div className="section-subtitle">{subtitle}</div>}
      </div>
    </div>
    {children}
  </div>
);


export default Section;