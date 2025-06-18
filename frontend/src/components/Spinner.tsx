export default function Spinner() {
  return (
    <div className="flex items-center justify-center space-x-2 animate-pulse">
      <div className="w-4 h-4 rounded-full bg-blue-600"></div>
      <div className="w-4 h-4 rounded-full bg-blue-600"></div>
      <div className="w-4 h-4 rounded-full bg-blue-600"></div>
    </div>
  );
}
