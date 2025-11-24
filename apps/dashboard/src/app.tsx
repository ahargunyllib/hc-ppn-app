import Layout from "./shared/components/layout";

export default function App() {
  return (
    <Layout>
      <Example />
    </Layout>
  );
}
function Example() {
  return (
    <div className="container mx-auto flex flex-col gap-2 py-20">
      <h1 className="font-bold text-3xl">Hello, Dashboard!</h1>
    </div>
  );
}
