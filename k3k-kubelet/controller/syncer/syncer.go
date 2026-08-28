package syncer

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/k3k/k3k-kubelet/translate"
)

type SyncerContext struct {
	ClusterName      string
	ClusterNamespace string
	VirtualClient    client.Client
	HostClient       client.Client
	Translator       translate.ToHostTranslator
	LabelSelector    labels.Selector
}

func GenerateLabelSelector(selector map[string]string, filter string) (labels.Selector, error) {
	mergedSelector := labels.Everything()
	labelSelector := labels.SelectorFromSet(selector)

	filterSelector, err := metav1.ParseToLabelSelector(filter)
	if err != nil {
		return nil, fmt.Errorf("error parsing filter selector %s: %w", filter, err)
	}

	filteredLabelSelector, err := metav1.LabelSelectorAsSelector(filterSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to convert metav1.labelSelector to labels.Selector %w", err)
	}

	labelSelectorRequirements, _ := labelSelector.Requirements()
	filteredSelectorRequirements, _ := filteredLabelSelector.Requirements()

	// merge the two requirements
	return mergedSelector.Add(append(labelSelectorRequirements, filteredSelectorRequirements...)...), nil
}
