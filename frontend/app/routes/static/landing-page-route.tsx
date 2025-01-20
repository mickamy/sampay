import { Link as LinkIcon, Share2, User } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import Header from "~/components/header";
import { Button } from "~/components/ui/button";
import FeatureCard from "~/routes/static/components/fature-card";
import { FAQ } from "./components/faq";

export default function LandingPage() {
  const { t } = useTranslation();
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-100 to-gray-200">
      <Header isLoggedIn={false} />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <section className="text-center mb-16">
          <h1 className="text-4xl md:text-5xl font-extrabold text-gray-900 mb-4 whitespace-pre-line tracking-wider leading-normal md:leading-normal">
            {t("lp.tagline")}
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            {t("lp.tagline_description_1")}
            <br />
            {t("lp.tagline_description_2")}
          </p>
          <Link to="/sign-up">
            <Button size="lg" className="mr-4">
              {t("lp.get_started")}
            </Button>
          </Link>
        </section>

        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <FeatureCard
            icon={<User className="h-12 w-12 text-blue-500" />}
            title={t("lp.feature.1.title")}
            description={t("lp.feature.1.description")}
          />
          <FeatureCard
            icon={<LinkIcon className="h-12 w-12 text-green-500" />}
            title={t("lp.feature.2.title")}
            description={t("lp.feature.2.description")}
          />
          <FeatureCard
            icon={<Share2 className="h-12 w-12 text-purple-500" />}
            title={t("lp.feature.3.title")}
            description={t("lp.feature.3.description")}
          />
        </div>

        <FAQ />
      </main>

      <footer className="bg-gray-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 grid md:grid-cols-4 gap-8">
          <div>
            <h3 className="text-lg font-semibold mb-4">サポート</h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="mailto:contact@sampay.link"
                  className="hover:text-gray-300"
                >
                  {t("lp.contact")}
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-lg font-semibold mb-4">{t("lp.legal")}</h3>
            <ul className="space-y-2">
              <li>
                <a href="/terms" className="hover:text-gray-300">
                  {t("lp.terms")}
                </a>
              </li>
              <li>
                <a href="/privacy" className="hover:text-gray-300">
                  {t("lp.privacy")}
                </a>
              </li>
            </ul>
          </div>
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-8 pt-8 border-t border-gray-700 text-center">
          <p>&copy; 2025 Sampay. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}
